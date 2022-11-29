package adapter

import (
	"keeper/app/db"
	"keeper/app/db/adapter/mongo"
	"keeper/app/db/adapter/mysql"

	mongoAnalyser "keeper/app/plugins/pluginMongdb/backend"
	mysqlAnalyser "keeper/app/plugins/pluginMysql/backend"
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
