package sideQuests

import (
	"keeper/app/code"
	"keeper/app/modules"
	"keeper/app/pkg/standard"
	"keeper/app/schema"
	"keeper/app/tools"

	"github.com/mitchellh/mapstructure"
)

var databaseLast code.UnixTime

var analysedStructure *schema.DatabaseInfo

var analysedTime code.UnixTime = 0
var loadingModel bool
var statusCounter int

func getStatusCounter() int {
	statusCounter += 1
	return statusCounter
}

type DatabaseConnectionHandlers struct {
	Ch chan *modules.EchoMessage
}

func NewDatabaseConnectionHandlers(ch chan *modules.EchoMessage) *DatabaseConnectionHandlers {
	return &DatabaseConnectionHandlers{ch}
}

func (msg *DatabaseConnectionHandlers) Connect(args map[string]interface{}, structure *schema.DatabaseInfo) {
	connection := args["connection"].(map[string]interface{})
	database := args["database"].(string)

	databaseLast = tools.NewUnixTime()
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
		driver, err = NewMysqlDriver(connection)
		if err != nil {
			msg.setStatus(&StatusMessage{Name: "err", Message: err.Error()})
			return
		}
	case standard.MONGOALIAS:
		driver, err = NewMongoDriver(connection)
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
	statusString := tools.ToJsonStr(status)
	if serverlastStatus != statusString {
		msg.Ch <- &modules.EchoMessage{
			MsgType: "status",
			Payload: map[string]interface{}{
				"status":  status,
				"counter": getStatusCounter(),
			},
		}
		serverlastStatus = statusString
	}
}

func (msg *DatabaseConnectionHandlers) readVersion(pool standard.SqlStandard) error {
	version, err := pool.GetVersion()
	if err != nil {
		return err
	}

	msg.Ch <- &modules.EchoMessage{
		Payload: version,
		MsgType: "version",
		Dialect: pool.Dialect(),
	}

	return nil
}

func (msg *DatabaseConnectionHandlers) handleFullRefresh(pool standard.SqlStandard, strings ...string) {
	loadingModel = true
	msg.setStatusName("loadStructure")

	analysedTime = tools.NewUnixTime()

	tables, err := pool.Tables(strings...)
	if err == nil {
		msg.Ch <- &modules.EchoMessage{MsgType: "structure", Payload: tables}
	}

	msg.Ch <- &modules.EchoMessage{MsgType: "structureTime", Payload: analysedTime}
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
	analysedTime = tools.NewUnixTime()

	if forceSend || tables != nil {
		msg.Ch <- &modules.EchoMessage{
			MsgType: "structure",
			Payload: map[string]interface{}{
				"collections": tables,
				"engine":      pool.Dialect(),
			},
		}
	}

	msg.Ch <- &modules.EchoMessage{
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

	msg.Ch <- &modules.EchoMessage{
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
	databaseLast = tools.NewUnixTime()
}
