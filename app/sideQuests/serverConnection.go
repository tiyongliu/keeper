package sideQuests

import (
	"keeper/app/code"
	"keeper/app/modules"
	"keeper/app/pkg/logger"
	"keeper/app/pkg/serializer"
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

type ServerConnectionHandlers struct {
	Mysql standard.SqlStandard
	Mongo standard.SqlStandard
	Ch    chan *modules.EchoMessage
}

func NewServerConnectionHandlers(ch chan *modules.EchoMessage) *ServerConnectionHandlers {
	//childProcessChecker(ch)
	setInterval(func() {
		close(ch)
	})

	return &ServerConnectionHandlers{
		Ch: ch,
	}
}

/*
定时器 是当你想要在未来某一刻执行一次时使用的
打点器 则是当你想要在固定的时间间隔重复执行准备的。这里是一个打点器的例子，它将定时的执行，直到我们将它停止。
*/
func setInterval(fn func()) {
	ticker := time.NewTicker(time.Minute)
	go func() {
		for range ticker.C {
			nowTime := time.Now().Unix()
			if code.UnixTime(nowTime)-serverlastPing > code.UnixTime(120*1000) {
				fn()
				ticker.Stop()
			}
		}
	}()

}

func (msg *ServerConnectionHandlers) Connect(connection map[string]interface{}) {
	msg.setStatusName("pending")
	serverlastPing = code.UnixTime(time.Now().Unix())
	//TODO request to dbEngineDriver
	//utility.RequireEngineDriver(connection)

	simpleSettingMysql := &modules.SimpleSettingMysql{}
	err := mapstructure.Decode(connection, simpleSettingMysql)
	if err != nil {
		return
	}

	//TODO connectUtility, 可以传递一个func 因为返回值都是一样的，在func内部进行处理
	var driver standard.SqlStandard
	switch connection["engine"].(string) {
	case code.MYSQLALIAS:
		driver, err = NewMysqlDriver(connection)
		if err != nil {
			logger.Infof("err: %v", err)
			msg.setStatus(&StatusMessage{
				Name:    "error",
				Message: err.Error(),
			})
			msg.errorExit()
			return
		}
		msg.Mysql = driver
	case code.MONGOALIAS:
		driver, err = NewMongoDriver(connection)
		if err != nil {
			msg.setStatus(&StatusMessage{
				Name:    "error",
				Message: err.Error(),
			})
			msg.errorExit()
			return
		}
		msg.Mongo = driver
	}

	if err := msg.readVersion(driver); err != nil {
		msg.setStatus(&StatusMessage{
			Name:    "error",
			Message: err.Error(),
		})
		msg.errorExit()
		return
	}

	if err := msg.handleRefresh(driver); err != nil {
		msg.setStatus(&StatusMessage{
			Name:    "error",
			Message: err.Error(),
		})
		msg.errorExit()
		return
	}

	msg.setStatusName("ok")
}

func (msg *ServerConnectionHandlers) Ping() code.UnixTime {
	return code.UnixTime(time.Now().Unix())
}

func (msg *ServerConnectionHandlers) CreateDatabase() {

}

func (msg *ServerConnectionHandlers) setStatusName(name string, message ...string) {
	if len(message) == 0 {
		msg.setStatus(&StatusMessage{name, ""})
	} else {
		msg.setStatus(&StatusMessage{name, message[0]})
	}
}

func (msg *ServerConnectionHandlers) setStatus(status *StatusMessage) {
	statusString := tools.ToJsonStr(status)
	if serverlastStatus != statusString {
		msg.Ch <- &modules.EchoMessage{
			MsgType: "status",
			Payload: status,
		}
		serverlastStatus = statusString
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

func (msg *ServerConnectionHandlers) readVersion(pool standard.SqlStandard) error {
	version, err := pool.GetVersion()
	if err != nil {
		return err
	}

	msg.Ch <- &modules.EchoMessage{
		Payload: version,
		MsgType: "version",
	}

	return nil
}

func (msg *ServerConnectionHandlers) handleRefresh(pool standard.SqlStandard) error {
	databases, err := pool.ListDatabases()
	msg.setStatusName("ok")
	databasesString := tools.ToJsonStr(databases)
	if err != nil {
		return err
	}

	logger.Infof("chan handleRefresh databases -<: %s", tools.ToJsonStr(databases))

	if serverlastDatabases != databasesString {
		//TODO send
		msg.Ch <- &modules.EchoMessage{
			Payload: databases,
			MsgType: "databases",
			Dialect: pool.Dialect(),
		}
		serverlastDatabases = databasesString
	}

	return nil
}

func (msg *ServerConnectionHandlers) errorExit() {
	defer close(msg.Ch)
	timer := time.AfterFunc(1*time.Second, func() {
		msg.Ch <- &modules.EchoMessage{
			Payload: serializer.StatusCodeFailed,
			MsgType: "exit",
		}
	})

	defer timer.Stop()

}
