package sideQuests

import (
	"keeper/app/code"
	"keeper/app/modules"
	"keeper/app/pkg/logger"
	"keeper/app/pkg/standard"
	"keeper/app/plugins/pluginMongdb"
	"keeper/app/plugins/pluginMysql"
	"keeper/app/tools"
	"keeper/app/utility"
	"time"

	"github.com/mitchellh/mapstructure"
)

var serverlastStatus string
var serverlastPing code.UnixTime

var serverlastDatabases string

type StatusMessage struct {
	Name    string `json:"name"`
	Message string `json:"message"`
}

type ServerConnection struct {
	SqlDriver standard.SqlStandard
}

func NewServerConnection(ch chan *modules.EchoMessage) *ServerConnection {
	setInterval(func() {
		logger.Info("Server connection not alive, exiting")
		ch <- &modules.EchoMessage{
			MsgType: "exit",
		}
	})

	return &ServerConnection{}
}

/*
定时器 是当你想要在未来某一刻执行一次时使用的
打点器 则是当你想要在固定的时间间隔重复执行准备的。这里是一个打点器的例子，它将定时的执行，直到我们将它停止。
*/
func setInterval(fn func()) {
	ticker := time.NewTicker(time.Minute)
	go func(ticker *time.Ticker) {
		for range ticker.C {
			nowTime := tools.NewUnixTime()
			if nowTime-serverlastPing > code.UnixTime(120*1000) {
				fn()
				ticker.Stop()
			}
		}
	}(ticker)
}

func (msg *ServerConnection) Connect(ch chan *modules.EchoMessage, connection map[string]interface{}) {
	msg.setStatus(ch, "pending")

	serverlastPing = tools.NewUnixTime()

	//TODO connectUtility, 可以传递一个func 因为返回值都是一样的，在func内部进行处理

	sqlDriver, err := GetSqlDriver(connection)
	if err != nil {
		msg.setStatus(ch, "error", err.Error())
		return
	}

	msg.SqlDriver = sqlDriver

	if err := msg.readVersion(ch, sqlDriver); err != nil {
		msg.setStatus(ch, "error", err.Error())
		logger.Infof("readVersion err: [%v]", err)
		//msg.errorExit()
		return
	}

	if err := msg.handleRefresh(ch, sqlDriver); err != nil {
		msg.setStatus(ch, "error", err.Error())
		logger.Infof("handleRefresh err: [%v]", err)
		return
	}

	msg.setStatus(ch, "ok")
	//ticker := time.NewTicker(time.Second)
	//go func(ticker *time.Ticker) {
	//	for range ticker.C {
	//		if err := sqlDriver.Ping(); err != nil {
	//			ch <- &modules.EchoMessage{
	//				Payload: nil,
	//				MsgType: "exit",
	//				Dialect: sqlDriver.Dialect(),
	//			}
	//			ticker.Stop()
	//		}
	//	}
	//}(ticker)
}

func GetSqlDriver(connection map[string]interface{}) (driver standard.SqlStandard, err error) {
	switch connection["engine"].(string) {
	case standard.MYSQLALIAS:
		driver, err = NewMysqlDriver(connection)
		if err != nil {
			return nil, err
		}
	case standard.MONGOALIAS:
		driver, err = NewMongoDriver(connection)
		if err != nil {
			return nil, err
		}
	}

	return driver, nil
}

func (msg *ServerConnection) NewTime() {
	serverlastPing = tools.NewUnixTime()
}

func (msg *ServerConnection) Ping(connection map[string]interface{}) *modules.OpenedStatus {
	if msg.SqlDriver == nil {
		sqlDriver, err := GetSqlDriver(connection)
		if err != nil {
			return &modules.OpenedStatus{
				Name:    "error",
				Message: err.Error(),
			}
		}

		if sqlDriver.Ping() != nil {
			return &modules.OpenedStatus{
				Name:    "error",
				Message: err.Error(),
			}
		}

		return &modules.OpenedStatus{
			Name:    "ok",
			Message: "",
		}
	}

	return nil
}

func (msg *ServerConnection) CreateDatabase() {

}

func (msg *ServerConnection) setStatus(ch chan *modules.EchoMessage, name string, message ...string) {
	status := map[string]string{"name": name}
	if len(message) > 0 {
		status["message"] = message[0]
	}

	statusString := tools.ToJsonStr(status)
	if serverlastStatus != statusString {
		ch <- &modules.EchoMessage{
			MsgType: "status",
			Payload: status,
		}
		serverlastStatus = statusString
	}

	if name == "error" {
		ch <- &modules.EchoMessage{
			MsgType: "exit",
		}
	}
}

func NewMysqlDriver(connection map[string]interface{}) (standard.SqlStandard, error) {
	storedConnection := connectUtility(connection)
	simpleSettingMysql := &modules.SimpleSettingMysql{}
	err := mapstructure.Decode(storedConnection, simpleSettingMysql)
	if err != nil {
		return nil, err
	}

	pool, err := pluginMysql.NewSimpleMysqlPool(simpleSettingMysql)
	if err != nil {
		return nil, err
	}

	return pool, nil
}

func NewMongoDriver(connection map[string]interface{}) (standard.SqlStandard, error) {
	storedConnection := connectUtility(connection)
	pool, err := pluginMongdb.NewSimpleMongoDBPool(&modules.SimpleSettingMongoDB{
		Host: storedConnection["host"],
		Port: storedConnection["port"],
	})

	if err != nil {
		return nil, err
	}

	return pool, nil
}

func connectUtility(connection map[string]interface{}) map[string]string {
	return utility.DecryptConnection(tools.TransformStringMap(connection))
}

func (msg *ServerConnection) readVersion(ch chan *modules.EchoMessage, pool standard.SqlStandard) error {
	version, err := pool.GetVersion()
	if err != nil {
		return err
	}

	ch <- &modules.EchoMessage{
		Payload: version,
		MsgType: "version",
	}

	return nil
}

func (msg *ServerConnection) handleRefresh(ch chan *modules.EchoMessage, pool standard.SqlStandard) error {
	databases, err := pool.ListDatabases()
	msg.setStatus(ch, "ok")
	databasesString := tools.ToJsonStr(databases)
	if err != nil {
		return err
	}

	if serverlastDatabases != databasesString {
		ch <- &modules.EchoMessage{
			Payload: databases,
			MsgType: "databases",
			Dialect: pool.Dialect(),
		}
		serverlastDatabases = databasesString
	}

	return nil
}

func (msg *ServerConnection) errorExit() {
	//defer close(msg.Ch)
	//timer := time.AfterFunc(1*time.Second, func() {
	//	msg.Ch <- &modules.EchoMessage{
	//		Payload: serializer.StatusCodeFailed,
	//		MsgType: "exit",
	//	}
	//})
	//
	//defer timer.Stop()
}
