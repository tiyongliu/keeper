package backend

import (
	"fmt"
	"keeper/app/pkg/logger"
	"keeper/app/plugins/modules"
	"keeper/app/plugins/pluginMysql"
	"keeper/app/utility"
	"path"
	"testing"
)

func TestRun(t *testing.T) {
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
	//NewAnalyser(pool, "yami_shops").RunAnalysis()

	analyser := NewAnalyser(pool, "yami_shops")

	_runAnalysis := analyser.RunAnalysis()

	resp := analyser.DatabaseAnalyser.AddEngineField(_runAnalysis)
	dir := utility.DataDir()

	if err = utility.WriteFileAllPool(utility.NewJsonLinesDatabase(path.Join(dir, "yami_shops.jsonl")).Filename, []map[string]interface{}{resp}); err != nil {
		logger.Errorf("err: %v", err)
	}
}
