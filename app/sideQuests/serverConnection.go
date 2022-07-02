package sideQuests

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"keeper/app/code"
	"keeper/app/modules"
	"keeper/app/pkg/logger"
	"keeper/app/pkg/standard"
	plugin_mondb "keeper/app/plugins/plugin-mondb"
	plugin_mysql "keeper/app/plugins/plugin-mysql"
	"keeper/app/tools"
	"keeper/app/utility"
	"time"
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
	return &MessageDriverHandlers{
		Ch: ch,
	}
}

/*
定时器 是当你想要在未来某一刻执行一次时使用的
打点器 则是当你想要在固定的时间间隔重复执行准备的。这里是一个打点器的例子，它将定时的执行，直到我们将它停止。
*/
func (msg *MessageDriverHandlers) Start() {
	ticker := time.NewTicker(time.Minute)
	go func() {
		for range ticker.C {
			nowTime := time.Now().Unix()
			if code.UnixTime(nowTime)-lastPing > code.UnixTime(120*1000) {
				logger.Info("Server connection not alive, exiting")
				//todo process.exit(0);
				ticker.Stop()
			}
		}
	}()

}

func (msg *MessageDriverHandlers) Connect(connection map[string]interface{}) {
	defer close(msg.Ch)
	msg.setStatusName("pending")
	lastPing = code.UnixTime(time.Now().Unix())

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
			return
		}
		msg.Mysql = driver
	case code.MONGOALIAS:
		driver, err = NewMongoDriver(connection)
		if err != nil {
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
		fmt.Println(tools.ToJsonStr(version))
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

	if lastDatabases != databasesString {
		//TODO send
		msg.Ch <- &modules.EchoMessage{
			Payload: &modules.DriverPayload{
				Name:        pool.Dialect(),
				StandardRes: databases,
			},
			MsgType: "databases",
		}
		lastDatabases = databasesString
	}

	return nil
}

func (msg *MessageDriverHandlers) errorExit() {
	defer close(msg.Ch)
	timer := time.AfterFunc(1*time.Second, func() {
		msg.Ch <- &modules.EchoMessage{
			Payload: 1,
			MsgType: "exit",
		}
	})

	defer timer.Stop()

}
