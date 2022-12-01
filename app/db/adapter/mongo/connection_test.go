package mongo

import (
	"fmt"
	"keeper/app/pkg/logger"
	"keeper/app/utility"
	"testing"
)

func TestParseSetting(t *testing.T) {
	setting, err := ParseSetting(map[string]interface{}{
		//"username": "root",
		//"password": "123456",
		"host": "localhost",
		//"port":     "27017",
	})
	if err != nil {
		logger.Infof("err: %v", err)
		return
	}

	fmt.Println(setting.Port)
	logger.Infof("setting: %s", utility.ToJsonStr(setting))
	logger.Infof("%s", setting.String())
}
