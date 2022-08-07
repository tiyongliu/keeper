package utility

import (
	"errors"
	"github.com/mitchellh/mapstructure"
	"keeper/app/pkg/standard"
	"keeper/app/plugins/modules"
	"keeper/app/plugins/pluginMongdb"
	"keeper/app/plugins/pluginMysql"
	"sync"
)

var storedConnection sync.Map

func CreateEngineDriver(connection map[string]interface{}) (driver standard.SqlStandard, err error) {
	switch connection["engine"].(string) {
	case standard.MYSQLALIAS:
		driver, err = newMysqlDriver(connection)
		if err != nil {
			return nil, err
		}
	case standard.MONGOALIAS:
		driver, err = newMongoDriver(connection)
		if err != nil {
			return nil, err
		}
	}

	return driver, nil
}

func newMysqlDriver(connection map[string]interface{}) (standard.SqlStandard, error) {
	loadedWithDb := connectUtility(connection)
	simpleSettingMysql := &modules.SimpleSettingMysql{}
	err := mapstructure.Decode(loadedWithDb, simpleSettingMysql)
	if err != nil {
		return nil, err
	}

	pool, err := pluginMysql.NewSimpleMysqlPool(simpleSettingMysql)
	if err != nil {
		return nil, err
	}

	return pool, nil
}

func newMongoDriver(connection map[string]interface{}) (standard.SqlStandard, error) {
	loadedWithDb := connectUtility(connection)
	pool, err := pluginMongdb.NewSimpleMongoDBPool(&modules.SimpleSettingMongoDB{
		Host: loadedWithDb["host"],
		Port: loadedWithDb["port"],
	})

	if err != nil {
		return nil, err
	}

	return pool, nil
}

func connectUtility(connection map[string]interface{}) map[string]string {
	return DecryptConnection(TransformStringMap(connection))
}

func SetDriverPool(conid string, driver standard.SqlStandard) error {
	if driver == nil {
		return errors.New("invalid memory address or nil pointer dereference")
	}
	storedConnection.Store(conid, driver)
	return nil
}

func GetDriverPool(conid string) (driver standard.SqlStandard, err error) {
	load, ok := storedConnection.Load(conid)
	if !ok {
		return nil, errors.New("invalid memory address or nil pointer dereference")
	}
	sqlStandard := load.(standard.SqlStandard)
	if sqlStandard == nil {
		return nil, errors.New("invalid memory address or nil pointer dereference")
	}

	return sqlStandard, nil
}

func DeleteDriverPool(conid string) error {
	andDelete, ok := storedConnection.LoadAndDelete(conid)
	if ok {
		if andDelete != nil {
			return andDelete.(standard.SqlStandard).Close()
		}
	}
	return nil
}