package backend

import (
	"github.com/samber/lo"
	"keeper/app/db"
	"keeper/app/db/adapter/mongo"
	"keeper/app/db/standard/modules"
	"keeper/app/plugins"
)

type Analyser struct {
	Driver           db.Session
	DatabaseAnalyser *plugins.DatabaseAnalyser
	DatabaseName     string
}

func NewAnalyser(driver db.Session, database string) *Analyser {
	return &Analyser{
		Driver:           driver,
		DatabaseName:     database,
		DatabaseAnalyser: plugins.NewDatabaseAnalyser(driver),
	}
}

func (da *Analyser) RunAnalysis() map[string]interface{} {
	driver, ok := da.Driver.(*mongo.Source)
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
