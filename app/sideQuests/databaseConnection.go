package sideQuests

import (
	"github.com/mitchellh/mapstructure"
	"keeper/app/code"
	"keeper/app/modules"
	"keeper/app/pkg/logger"
	"keeper/app/pkg/standard"
	"time"
)

var databaseLastPing code.UnixTime

type DatabaseConnectionHandlers struct {
	Ch chan interface{}
}

func NewDatabaseConnectionHandlers() *DatabaseConnectionHandlers {
	return &DatabaseConnectionHandlers{}
}

func (msg *DatabaseConnectionHandlers) Connect(connection map[string]interface{}, structure string) {
	databaseLastPing = code.UnixTime(time.Now().Unix())

	if structure == "" {

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
}
