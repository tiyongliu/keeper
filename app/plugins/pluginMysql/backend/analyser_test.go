package backend

import (
	"fmt"
	"keeper/app/db/drivers"
	"keeper/app/pkg/logger"
	"keeper/app/utility"
	"path"
	"testing"
)

func TestRun(t *testing.T) {
	driver, err := drivers.NewCompatDriver().Open(map[string]interface{}{
		"username": "root",
		"password": "123456",
		"port":     "3306",
		"database": "",
		"host":     "localhost",
	})

	if err != nil {
		fmt.Printf("err: %v \n", err)
		return
	}
	//NewAnalyser(pool, "yami_shops").RunAnalysis()

	analyser := NewAnalyser(driver, "yami_shops")

	_runAnalysis := analyser.RunAnalysis()

	resp := analyser.DatabaseAnalyser.AddEngineField(_runAnalysis)
	dir := utility.DataDir()

	if err = utility.WriteFileAllPool(utility.NewJsonLinesDatabase(path.Join(dir, "yami_shops.jsonl")).Filename, []map[string]interface{}{resp}); err != nil {
		logger.Errorf("err: %v", err)
	}
}
