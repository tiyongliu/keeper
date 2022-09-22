package adapter

import (
	"keeper/app/pkg/standard"
	mongoAnalyser "keeper/app/plugins/pluginMongdb/backend"
	mysqlAnalyser "keeper/app/plugins/pluginMysql/backend"
)

func AnalyseFull(driver standard.SqlStandard, database string) map[string]interface{} {
	switch driver.Dialect() {
	case standard.MONGOALIAS:
		analyser := mongoAnalyser.NewAnalyser(driver, database)
		return analyser.DatabaseAnalyser.AddEngineField(analyser.RunAnalysis())
	case standard.MYSQLALIAS:
		analyser := mysqlAnalyser.NewAnalyser(driver, database)
		return analyser.DatabaseAnalyser.AddEngineField(analyser.RunAnalysis())
	default:
		return nil
	}
}

func AnalyseIncremental() {

}

//查询表
func CreateDumper() {

}
