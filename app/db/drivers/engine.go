package drivers

import (
	"github.com/samber/lo"
	"keeper/app/db"
	"keeper/app/db/adapter/mongo"
	"keeper/app/db/adapter/mysql"
	"keeper/app/internal"
	"keeper/app/pkg/logger"
	"keeper/app/utility"
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
	return createSession(internal.DecryptConnection(loadConnection(storedConnection)))
}

func loadConnection(storedConnection map[string]interface{}) map[string]interface{} {
	if storedConnection[uuid] != nil && storedConnection[uuid].(string) != "" {
		loaded := internal.GetCore(storedConnection[uuid].(string), false)
		return lo.Assign(loaded, storedConnection)
	}

	return storedConnection
}

func createSession(setting map[string]interface{}) (db.Session, error) {
	if setting != nil && setting[driverName].(string) != "" {
		switch setting[driverName] {
		case mysql.Adapter:
			logger.Infof("setting = %s", utility.ToJsonStr(setting))
			parseSetting, err := mysql.ParseSetting(setting)
			if err != nil {
				logger.Errorf("setting parse failed %v", err)
				return nil, err
			}
			return mysql.Open(parseSetting)
		case mongo.Adapter:
			parseSetting, err := mongo.ParseSetting(setting)
			if err != nil {
				logger.Errorf("setting parse failed %v", err)
				return nil, err
			}
			return mongo.Open(parseSetting)
		}
	}
	return nil, db.ErrMissingDriverName
}
