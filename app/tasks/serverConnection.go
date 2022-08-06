package tasks

import (
	"github.com/mitchellh/mapstructure"
	"keeper/app/pkg/standard"
	"keeper/app/plugins/modules"
	"keeper/app/plugins/pluginMongdb"
	"keeper/app/plugins/pluginMysql"
	"keeper/app/utility"
)

type ServerConnection struct {
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
	return utility.DecryptConnection(utility.TransformStringMap(connection))
}
