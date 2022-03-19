package backend

import (
	"github.com/samber/lo"
	"keeper/app/pkg/standard"
	"keeper/app/plugins"
	"keeper/app/plugins/modules"
	"keeper/app/plugins/pluginMongdb"
)

type Analyser struct {
	Driver           standard.SqlStandard
	DatabaseAnalyser *plugins.DatabaseAnalyser
	DatabaseName     string
}

func NewAnalyser(driver standard.SqlStandard, database string) *Analyser {
	return &Analyser{
		Driver:           driver,
		DatabaseName:     database,
		DatabaseAnalyser: plugins.NewDatabaseAnalyser(driver),
	}
}

func (da *Analyser) RunAnalysis() map[string]interface{} {
	driver, ok := da.Driver.(*pluginMongdb.MongoDBDrivers)
	if !ok && driver == nil {
		return nil
	}
	collections, err := driver.Collections(da.DatabaseName)
	if err != nil {
		return nil
	}

	return da.DatabaseAnalyser.MergeAnalyseResult(map[string]interface{}{
		"collections": lo.Map(collections, func(x *modules.MongoDBCollection, i int) map[string]interface{} {
			return map[string]interface{}{"pureName": x.PureName}
		}),
	})
}
