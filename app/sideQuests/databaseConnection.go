package sideQuests

import (
	"keeper/app/code"
	"keeper/app/driver"
	"keeper/app/modules"
	"keeper/app/pkg/logger"
	"keeper/app/pkg/standard"
	"keeper/app/tools"
	"time"

	"github.com/mitchellh/mapstructure"
)

var databaseLastPing code.UnixTime

var analysedStructure *driver.Structure

var analysedTime code.UnixTime = 0

type DatabaseConnectionHandlers struct {
	Ch chan interface{}
}

func NewDatabaseConnectionHandlers() *DatabaseConnectionHandlers {
	return &DatabaseConnectionHandlers{}
}

func (msg *DatabaseConnectionHandlers) Connect(connection map[string]interface{}, structure *driver.Structure) {
	databaseLastPing = code.UnixTime(time.Now().Unix())
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
	case code.MYSQLALIAS:
		driver, err = NewMysqlDriver(connection)
		if err != nil {
			logger.Infof("err: %v", err)

			return
		}
	case code.MONGOALIAS:
		driver, err = NewMongoDriver(connection)
		if err != nil {

			return
		}
	}

	logger.Infof("driver info : %s", driver)

	if structure != nil {

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
			Payload: status,
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

func (msg *DatabaseConnectionHandlers) handleFullRefresh(pool standard.SqlStandard) {
	msg.setStatusName("loadStructure")

	analysedTime = code.UnixTime(time.Now().Unix())
	msg.setStatusName("ok")
}

func (msg *DatabaseConnectionHandlers) handleIncrementalRefresh(forceSend bool, pool standard.SqlStandard) {
	msg.setStatusName("checkStructure")
}
