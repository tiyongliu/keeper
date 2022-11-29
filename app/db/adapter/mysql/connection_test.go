package mysql

import (
	"keeper/app/pkg/logger"
	"testing"
)

func Test_ParseSetting(t *testing.T) {
	setting, err := ParseSetting(map[string]interface{}{
		"username": "root",
		"password": "123456",
		"port":     "3306",
		"database": "",
		"host":     "localhost",
	})

	if err != nil {
		logger.Infof("err: %v", err)
		return
	}

	logger.Infof("%s", setting.String())
}
