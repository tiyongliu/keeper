package adapter

import (
	"keeper/app/analyser/mongoAnalyser"
	"keeper/app/analyser/mysqlAnalyser"
	"keeper/app/db"
	"keeper/app/db/adapter/mongo"
	"keeper/app/db/adapter/mysql"
)

func AnalyseFull(driver db.Session, database string) map[string]interface{} {
	switch driver.Dialect() {
	case mongo.Adapter:
		analyser := mongoAnalyser.NewAnalyser(driver, database)
		return analyser.DatabaseAnalyser.AddEngineField(analyser.RunAnalysis())
	case mysql.Adapter:
		analyser := mysqlAnalyser.NewAnalyser(driver, database)
		return analyser.DatabaseAnalyser.AddEngineField(analyser.RunAnalysis())
	default:
		return nil
	}
}

func AnalyseIncremental() {

}

// 查询表
func CreateDumper() {

}
