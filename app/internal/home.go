package internal

import (
	"errors"
	"keeper/app/db"
	"keeper/app/db/adapter/mongo"
	"keeper/app/db/adapter/mysql"
	"keeper/app/pkg/logger"
	"keeper/app/utility"
	"path"
	"sync"
)

var storedConnection sync.Map

var JsonLinesDatabase *utility.JsonLinesDatabase

func init() {
	dir := utility.DataDir()
	JsonLinesDatabase = utility.NewJsonLinesDatabase(path.Join(dir, "connections.jsonl"))
}

func GetCore(conid string, mask bool) map[string]interface{} {
	if conid == "" {
		return nil
	}
	return JsonLinesDatabase.Get(conid)
}

func CreateEngineDriver(connection map[string]interface{}) (driver db.Session, err error) {
	logger.Infof("%s", utility.ToJsonStr(connection))
	utility.WithRecover(func() {
		switch connection["engine"].(string) {
		case mongo.Adapter:
		case mysql.Adapter:
		default:
			err = errors.New("invalid connection")
		}
	}, func(e error) {
		logger.Errorf("err: %v", e)
		err = e
	})
	return driver, err
}

func SetDriverPool(conid string, driver db.Session) error {
	if driver == nil {
		return errors.New("invalid memory address or nil pointer dereference")
	}
	storedConnection.Store(conid, driver)
	return nil
}

func GetStoragePool(conid string) (driver db.Session, err error) {
	load, ok := storedConnection.Load(conid)
	if !ok {
		return nil, errors.New("invalid memory address or nil pointer dereference")
	}
	sqlStandard := load.(db.Session)
	if sqlStandard == nil {
		return nil, errors.New("invalid memory address or nil pointer dereference")
	}

	return sqlStandard, nil
}

func DeleteStoragePool(conid string) error {
	andDelete, ok := storedConnection.LoadAndDelete(conid)
	if ok {
		if andDelete != nil {
			return andDelete.(db.Session).Close()
		}
	}
	return nil
}

func TargetStoragePool(conid string, connection map[string]interface{}) (driver db.Session, err error) {
	driver, err = GetStoragePool(conid)
	if driver == nil {
		driver, err = CreateEngineDriver(connection)
		err = SetDriverPool(conid, driver)
	}

	return
}

func CleanStoragePool() {
	storedConnection.Range(func(key, value any) bool {
		sqlStandard, ok := value.(db.Session)
		if sqlStandard != nil && ok {
			if err := sqlStandard.Close(); err != nil {
				logger.Infof("driver by conid: %s close failed: %v", key, err)
				return false
			}
		}
		return true
	})
}
