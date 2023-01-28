package adapter

import (
	"errors"
	"github.com/samber/lo"
	"keeper/app/db"
	"keeper/app/db/adapter/mongo"
	"keeper/app/db/adapter/mysql"
	"keeper/app/internal"
	"keeper/app/pkg/logger"
	"keeper/app/utility"
	"path/filepath"
)

const (
	driverName = `engine`
	uuid       = `_id`
)

type sessionWithContext struct {
}

type Driver interface {
	//visit simple drive connection
	Open(map[string]interface{}) (db.Session, error)
}

func NewCompatDriver() Driver {
	return &sessionWithContext{}
}

func (s *sessionWithContext) Open(storedConnection map[string]interface{}) (db.Session, error) {
	if !utility.IsExist(filepath.Join(utility.DataDir(), "connections.jsonl")) {
		return nil, errors.New("connections file missing")
	}
	return createSession(internal.DecryptConnection(loadConnection(storedConnection)))
}

func loadConnection(storedConnection map[string]interface{}) map[string]interface{} {
	if storedConnection[uuid] != nil && storedConnection[uuid].(string) != "" {
		loaded := internal.GetCore(storedConnection[uuid].(string), false)
		return lo.Assign(loaded, storedConnection)
	}

	return storedConnection
}

func createSession(storedConnection map[string]interface{}) (db.Session, error) {
	if storedConnection != nil && storedConnection[driverName].(string) != "" {
		switch storedConnection[driverName] {
		case mysql.Adapter:
			parseSetting, err := mysql.ParseSetting(storedConnection)
			if err != nil {
				logger.Errorf("setting parse failed %v", err)
				return nil, err
			}
			return mysql.Open(parseSetting)
		case mongo.Adapter:
			parseSetting, err := mongo.ParseSetting(storedConnection)
			if err != nil {
				logger.Errorf("setting parse failed %v", err)
				return nil, err
			}
			return mongo.Open(parseSetting)
		}
	}
	return nil, db.ErrMissingDriverName
}
