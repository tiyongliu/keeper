package mongo

import (
	"fmt"
	"keeper/app/db"
	"keeper/app/pkg/logger"
	"keeper/app/utility"
	"testing"
)

func getDevice(scope func(session db.Session)) db.Session {
	setting := &ConnectionURL{
		User:     "",
		Password: "",
		Host:     "localhost",
		Database: "",
		Port:     "27017",
		Options:  nil,
	}

	fmt.Println(setting.String())
	open, err := Open(setting)
	if err != nil {
		logger.Infof(Adapter+"connect failed: %v", err)
		panic(err)
	}
	scope(open)
	defer open.Close()

	return open
}

func TestDialect(t *testing.T) {
	getDevice(func(session db.Session) {
		fmt.Println(session.Dialect())
		fmt.Println(session.Ping())
	})

}

func TestGetVersion(t *testing.T) {
	getDevice(func(session db.Session) {
		version, err := session.Version()
		if err != nil {
			logger.Errorf("%v", err)
			return
		}
		logger.Infof("%s", utility.ToJsonStr(version))
	})

}
