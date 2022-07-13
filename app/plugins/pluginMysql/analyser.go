package pluginMysql

import (
	"github.com/samber/lo"
	"keeper/app/modules"
	"keeper/app/pkg/logger"
	"keeper/app/tools"
)

type Analyser struct {
	driver *MysqlDrivers
}

func NewAnalyser() *Analyser {
	return &Analyser{}
}

func (a *Analyser) Build() {
	databaseName := "yami_shops"
	tableName := "tz_user"
	tables, err := a.driver.Tables("yami_shops", "tz_user")
	if err != nil {
		return
	}

	columns, err := a.driver.Columns(databaseName, tableName)
	if err != nil {
		return
	}

	for _, table := range tables.([]*modules.MysqlTableSchema) {
		logger.Info("table info: %s", tools.ToJsonStr(table))
		lo.Filter[*modules.MysqlTableColumn](columns.([]*modules.MysqlTableColumn), func(col *modules.MysqlTableColumn, i int) bool {
			return table.PureName == col.PureName
		})
	}
}
