package sideQuests

import (
	"keeper/app/pkg/containers"
	"keeper/app/pkg/standard"
	"keeper/app/plugins/modules"
	"keeper/app/schema"
	"keeper/app/tasks"
	"keeper/app/utility"

	"github.com/mitchellh/mapstructure"
)

var databaseLast utility.UnixTime

var analysedStructure *schema.DatabaseInfo

var analysedTime utility.UnixTime = 0
var loadingModel bool
var statusCounter int

func getStatusCounter() int {
	statusCounter += 1
	return statusCounter
}

type DatabaseConnectionHandlers struct {
	Ch chan *containers.EchoMessage
}

func NewDatabaseConnectionHandlers(ch chan *containers.EchoMessage) *DatabaseConnectionHandlers {
	return &DatabaseConnectionHandlers{ch}
}

func (msg *DatabaseConnectionHandlers) Connect(args map[string]interface{}, structure *schema.DatabaseInfo) {
	connection := args["connection"].(map[string]interface{})
	database := args["database"].(string)

	databaseLast = utility.NewUnixTime()
	if structure == nil {
		msg.setStatusName("pending")
	}

	simpleSettingMysql := &modules.SimpleSettingMysql{}
	err := mapstructure.Decode(connection, simpleSettingMysql)
	if err != nil {
		return
	}

	var driver standard.SqlStandard
	switch connection["engine"].(string) {
	case standard.MYSQLALIAS:
		driver, err = tasks.NewMysqlDriver(connection)
		if err != nil {
			msg.setStatus(&StatusMessage{Name: "err", Message: err.Error()})
			return
		}
	case standard.MONGOALIAS:
		driver, err = tasks.NewMongoDriver(connection)
		if err != nil {
			msg.setStatus(&StatusMessage{Name: "err", Message: err.Error()})
			return
		}
	}

	if structure != nil {
		msg.handleIncrementalRefresh(true, driver, database)
	} else {
		msg.handleFullRefresh(driver, database)
	}

}

func (msg *DatabaseConnectionHandlers) setStatusName(name string, message ...string) {
	if len(message) == 0 {
		msg.setStatus(&StatusMessage{name, ""})
	} else {
		msg.setStatus(&StatusMessage{name, message[0]})
	}
}

func (msg *DatabaseConnectionHandlers) setStatus(status *StatusMessage) {
	statusString := utility.ToJsonStr(status)
	if ServerLastStatus != statusString {
		msg.Ch <- &containers.EchoMessage{
			MsgType: "status",
			Payload: map[string]interface{}{
				"status":  status,
				"counter": getStatusCounter(),
			},
		}
		ServerLastStatus = statusString
	}
}

func (msg *DatabaseConnectionHandlers) readVersion(pool standard.SqlStandard) error {
	version, err := pool.GetVersion()
	if err != nil {
		return err
	}

	msg.Ch <- &containers.EchoMessage{
		Payload: version,
		MsgType: "version",
		Dialect: pool.Dialect(),
	}

	return nil
}

func (msg *DatabaseConnectionHandlers) handleFullRefresh(pool standard.SqlStandard, strings ...string) {
	loadingModel = true
	msg.setStatusName("loadStructure")

	analysedTime = utility.NewUnixTime()

	tables, err := pool.Tables(strings...)
	if err == nil {
		msg.Ch <- &containers.EchoMessage{MsgType: "structure", Payload: tables}
	}

	msg.Ch <- &containers.EchoMessage{MsgType: "structureTime", Payload: analysedTime}
	msg.setStatusName("ok")
	loadingModel = false
}

func (msg *DatabaseConnectionHandlers) handleIncrementalRefresh(forceSend bool, pool standard.SqlStandard, args ...string) {
	msg.setStatusName("checkStructure")

	tables, err := pool.Tables(args...)
	if err != nil {
		msg.setStatusName("loadStructure", err.Error())
		return
	}
	analysedTime = utility.NewUnixTime()

	if forceSend || tables != nil {
		msg.Ch <- &containers.EchoMessage{
			MsgType: "structure",
			Payload: map[string]interface{}{
				"collections": tables,
				"engine":      pool.Dialect(),
			},
		}
	}

	msg.Ch <- &containers.EchoMessage{
		MsgType: "structureTime",
		Payload: analysedTime,
	}

	msg.setStatusName("ok")
}

func (msg *DatabaseConnectionHandlers) ReadVersion(pool standard.SqlStandard) error {
	version, err := pool.GetVersion()
	if err != nil {
		return err
	}

	msg.Ch <- &containers.EchoMessage{
		Payload: version,
		MsgType: "version",
	}

	return nil
}

func (msg *DatabaseConnectionHandlers) QueryData() {

}

func (msg *DatabaseConnectionHandlers) SyncModel() {
	if loadingModel {
		return
	}
}

func (msg *DatabaseConnectionHandlers) Ping() {
	databaseLast = utility.NewUnixTime()
}
