package backend

import (
	"fmt"
	"keeper/app/plugins/modules"
	"keeper/app/plugins/pluginMysql"
	"testing"
)

func TestNewAnalyser(t *testing.T) {
	pool, err := pluginMysql.NewSimpleMysqlPool(&modules.SimpleSettingMysql{
		Host:     "localhost",
		Username: "root",
		Password: "123456",
		Port:     "3306",
	})

	if err != nil {
		fmt.Printf("err: %v \n", err)
		return
	}
	NewAnalyser(pool, "yami_shops").RunAnalysis()
}
