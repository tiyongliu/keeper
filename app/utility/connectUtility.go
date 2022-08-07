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

var DriverPoolMap sync.Map

func CreateEngineDriver(connection map[string]interface{}) (driver standard.SqlStandard, err error) {
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
	return DecryptConnection(TransformStringMap(connection))
}

func SetDriverPool(conid string, driver standard.SqlStandard) error {
	if driver == nil {
		return errors.New("invalid memory address or nil pointer dereference")
	}
	DriverPoolMap.Store(conid, driver)
	return nil
}

func GetDriverPool(conid string) (driver standard.SqlStandard, err error) {
	load, ok := DriverPoolMap.Load(conid)
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
	andDelete, ok := DriverPoolMap.LoadAndDelete(conid)
	if ok {
		if andDelete != nil {
			return andDelete.(standard.SqlStandard).Close()
		}
	}
	return nil
}
