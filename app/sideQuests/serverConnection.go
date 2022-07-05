package sideQuests

import (
	"keeper/app/code"
	"keeper/app/modules"
	"keeper/app/pkg/logger"
	"keeper/app/pkg/serializer"
	"keeper/app/pkg/standard"
	plugin_mondb "keeper/app/plugins/plugin-mondb"
	plugin_mysql "keeper/app/plugins/plugin-mysql"
	"keeper/app/tools"
	"keeper/app/utility"
	"time"

	"github.com/mitchellh/mapstructure"
)

var lastStatus string
var lastPing code.UnixTime
var lastDatabases string

type StatusMessage struct {
	Name    string `json:"name"`
	Message string `json:"message"`
}

type MessageDriverHandlers struct {
	Mysql standard.SqlStandard
	Mongo standard.SqlStandard
	Ch    chan *modules.EchoMessage
}

func NewMessageDriverHandlers(ch chan *modules.EchoMessage) *MessageDriverHandlers {
	//childProcessChecker(ch)
	setInterval(ch)
	return &MessageDriverHandlers{
		Ch: ch,
	}
}

/*
定时器 是当你想要在未来某一刻执行一次时使用的
打点器 则是当你想要在固定的时间间隔重复执行准备的。这里是一个打点器的例子，它将定时的执行，直到我们将它停止。
*/
func setInterval(ch chan *modules.EchoMessage) {
	ticker := time.NewTicker(time.Minute)
	go func() {
		for range ticker.C {
			nowTime := time.Now().Unix()
			if code.UnixTime(nowTime)-lastPing > code.UnixTime(120*1000) {
				logger.Info("Server connection not alive, exiting")
				ch <- &modules.EchoMessage{
					Payload: serializer.StatusCodeFailed,
					MsgType: "exit",
				}
				ticker.Stop()
			}
		}
	}()

}

func (msg *MessageDriverHandlers) Connect(connection map[string]interface{}) {
	msg.setStatusName("pending")
	lastPing = code.UnixTime(time.Now().Unix())
	logger.Info("1 info logger to =======================================")
	//TODO request to dbEngineDriver
	//utility.RequireEngineDriver(connection)

	simpleSettingMysql := &modules.SimpleSettingMysql{}
	err := mapstructure.Decode(connection, simpleSettingMysql)
	if err != nil {
		return
	}
	logger.Info("2 info logger to =======================================")
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
	logger.Info("3 info logger to =======================================")
	if err := msg.readVersion(driver); err != nil {
		msg.setStatus(&StatusMessage{
			Name:    "error",
			Message: err.Error(),
		})
		msg.errorExit()
		return
	}

	logger.Info("4 info logger to =======================================")
	if err := msg.handleRefresh(driver); err != nil {
		msg.setStatus(&StatusMessage{
			Name:    "error",
			Message: err.Error(),
		})
		msg.errorExit()
		return
	}

	logger.Info("5 info logger to =======================================")
	msg.setStatusName("ok")
}

func (msg *MessageDriverHandlers) Ping() code.UnixTime {
	return code.UnixTime(time.Now().Unix())
}

func (msg *MessageDriverHandlers) CreateDatabase() {

}

func (msg *MessageDriverHandlers) setStatusName(name string, message ...string) {
	if len(message) == 0 {
		msg.setStatus(&StatusMessage{name, ""})
	} else {
		msg.setStatus(&StatusMessage{name, message[0]})
	}
}

func (msg *MessageDriverHandlers) setStatus(status *StatusMessage) {
	statusString := tools.ToJsonStr(status)
	if lastStatus != statusString {
		msg.Ch <- &modules.EchoMessage{
			MsgType: "status",
			Payload: status,
		}
		lastStatus = statusString
	}
}

func NewMysqlDriver(connection map[string]interface{}) (standard.SqlStandard, error) {
	storedConnection := connectUtility(connection)
	simpleSettingMysql := &modules.SimpleSettingMysql{}
	err := mapstructure.Decode(storedConnection, simpleSettingMysql)
	if err != nil {
		return nil, err
	}

	pool, err := plugin_mysql.NewSimpleMysqlPool(simpleSettingMysql)
	if err != nil {
		return nil, err
	}

	return pool, nil
}

func NewMongoDriver(connection map[string]interface{}) (standard.SqlStandard, error) {
	storedConnection := connectUtility(connection)
	pool, err := plugin_mondb.NewSimpleMongoDBPool(&modules.SimpleSettingMongoDB{
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

func (msg *MessageDriverHandlers) readVersion(pool standard.SqlStandard) error {
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

func (msg *MessageDriverHandlers) handleRefresh(pool standard.SqlStandard) error {
	databases, err := pool.ListDatabases()
	msg.setStatusName("ok")
	databasesString := tools.ToJsonStr(databases)
	if err != nil {
		return err
	}

	logger.Infof("chan handleRefresh databases -<: %s", tools.ToJsonStr(databases))

	if lastDatabases != databasesString {
		//TODO send
		msg.Ch <- &modules.EchoMessage{
			Payload: databases,
			MsgType: "databases",
			Dialect: pool.Dialect(),
		}
		lastDatabases = databasesString
	}

	return nil
}

func (msg *MessageDriverHandlers) errorExit() {
	defer close(msg.Ch)
	timer := time.AfterFunc(1*time.Second, func() {
		msg.Ch <- &modules.EchoMessage{
			Payload: serializer.StatusCodeFailed,
			MsgType: "exit",
		}
	})

	defer timer.Stop()

}
