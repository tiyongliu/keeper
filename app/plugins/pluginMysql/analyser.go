package pluginMysql

import (
	"keeper/app/pkg/standard"
)

type Analyser struct {
	driver standard.SqlStandard
}

func NewAnalyser(driver standard.SqlStandard) *Analyser {
	return &Analyser{driver}
}

func (a *Analyser) Build() {

	//databaseName := "yami_shops"
	//tableName := "tz_user"
	//tables, err := a.schema.Tables(databaseName, tableName)
	//if err != nil {
	//	return
	//}

	//columns, err := a.schema.Columns(databaseName, tableName)
	//if err != nil {
	//	return
	//}
	//
	//var res []*modules.MysqlTableColumn
	//for _, table := range tables.([]*modules.MysqlTableSchema) {
	//	logger.Info("table info: %s", tools.ToJsonStr(table))
	//	res = lo.Filter[*modules.MysqlTableColumn](columns.([]*modules.MysqlTableColumn), func(col *modules.MysqlTableColumn, i int) bool {
	//		return table.PureName == col.PureName
	//	})
	//}

}

//func AnalyseIncremental(standard standard.SqlStandard, databaseName, tableName string) (interface{}, error) {
//	drivers, ok := standard.(*MysqlDrivers)
//	if !ok {
//		return nil, errors.New(serializer.TypeConversionError)
//	}
//
//	tables, err := drivers.Tables(databaseName, tableName)
//	if err != nil {
//		return nil, err
//	}
//
//	return tables, nil
//}

func (a *Analyser) IncrementalAnalysis(args ...string) {

}
