package plugins

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"keeper/app/code"
	"keeper/app/modules"
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
}

func NewMessageDriverHandlers() *MessageDriverHandlers {
	return &MessageDriverHandlers{}
}

func (msg *MessageDriverHandlers) Connect(connection map[string]interface{}) {
	setStatusName("pending")
	lastPing = code.UnixTime(time.Now().Unix())

	//TODO request to dbEngineDriver
	//utility.RequireEngineDriver(connection)

	simpleSettingMysql := &modules.SimpleSettingMysql{}
	err := mapstructure.Decode(connection, simpleSettingMysql)
	if err != nil {
		return
	}

	//TODO connectUtility, 可以传递一个func 因为返回值都是一样的，在func内部进行处理
	switch connection["engine"].(string) {
	case code.Mysql_alias:
		driver, err := NewMysqlDriver(connection)
		if err != nil {
			return
		}
		msg.Mysql = driver
	case code.Mongo_alias:
		driver, err := NewMongoDriver(connection)
		if err != nil {
			return
		}
		msg.Mongo = driver
	}

	/*
	  readVersion()
	  handleRefresh()
	*/

	//msg.SystemConnection, err = NewSimpleMysqlPool(simpleSettingMysql)
	//if err != nil {
	//	setStatus(&StatusMessage{"error", err.Error()})
	//}
}

func (msg *MessageDriverHandlers) Ping() code.UnixTime {
	return code.UnixTime(time.Now().Unix())
}

func (msg *MessageDriverHandlers) CreateDatabase() {

}

func setStatusName(name string, message ...string) {
	if len(message) == 0 {
		setStatus(&StatusMessage{name, ""})
	} else {
		setStatus(&StatusMessage{name, message[0]})
	}
}

func setStatus(status *StatusMessage) {
	statusString := tools.ToJsonStr(status)
	if lastStatus != statusString {
		//TODO send 消息

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

//TODO send
func readVersion(pool standard.SqlStandard) {
	version, err := pool.GetVersion()
	if err == nil {
		fmt.Println(tools.ToJsonStr(version))
	}
}

func handleRefresh(pool standard.SqlStandard) {
	databases, err := pool.ListDatabases()

	databasesString := tools.ToJsonStr(databases)
	if err != nil {
		setStatus(&StatusMessage{Name: "ok", Message: ""})
		return
	}

	if lastDatabases != databasesString {
		//TODO send
		lastDatabases = databasesString
	}

}
