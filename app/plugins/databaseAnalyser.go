package plugins

import (
	"fmt"
	"github.com/samber/lo"
	"keeper/app/db"
	"keeper/app/db/standard/modules"
	"strings"
)

type DatabaseAnalyser struct {
	Modifications []string
	Driver        db.Session
	Structure     interface{}
}

type AnalyserAdapter interface {
	AddEngineField()
}

var structureFields = []string{"tables", "collections", "views", "matviews", "functions", "procedures", "triggers"}

func NewDatabaseAnalyser(driver db.Session) *DatabaseAnalyser {
	return &DatabaseAnalyser{
		Modifications: nil,
		Driver:        driver,
	}
}

func (da *DatabaseAnalyser) CreateQuery(template string, typeFields []string) string {
	if da.Modifications == nil || len(da.Modifications) == 0 || typeFields == nil || len(typeFields) == 0 {
		return strings.Replace(template, "=OBJECT_ID_CONDITION", " is not null", -1)
	}

	return strings.Replace(template, "=OBJECT_ID_CONDITION", fmt.Sprintf(" in (%s)", strings.Join(typeFields, ",")), -1)
}

func (da *DatabaseAnalyser) AddEngineField(db map[string]interface{}) map[string]interface{} {
	if da.Driver == nil {
		return nil
	}

	for _, field := range structureFields {
		if db[field] == nil {
			continue
		}

		items, ok := db[field].([]map[string]interface{})
		if !ok || len(items) == 0 {
			continue
		}

		for _, item := range items {
			item["engine"] = da.Driver.Dialect()
		}
	}
	db["engine"] = da.Driver.Dialect()

	return db
}

func (da *DatabaseAnalyser) MergeAnalyseResult(newlyAnalysed map[string]interface{}) map[string]interface{} {
	return lo.Assign(CreateEmptyStructure(), newlyAnalysed)
}

func ExtractPrimaryKeys(table *modules.Table, pkColumns []*modules.PrimaryKey) map[string]interface{} {
	filtered := lo.Filter[*modules.PrimaryKey](pkColumns, func(x *modules.PrimaryKey, _ int) bool {
		return x.PureName == table.PureName
	})

	if len(filtered) == 0 {
		return nil
	}
	result := map[string]interface{}{"constraintType": "primaryKey"}
	if len(filtered) > 0 {
		result["constraintName"] = filtered[0].ConstraintName
		result["pureName"] = filtered[0].PureName
	}

	var keyColumns []map[string]string
	for _, key := range filtered {
		keyColumns = append(keyColumns, map[string]string{
			"columnName": key.ColumnName,
		})
	}
	result["columns"] = keyColumns
	return result
}

func ExtractForeignKeys(table *modules.Table, fkColumns []*modules.ForeignKeys) []map[string]interface{} {
	array := lo.Filter[*modules.ForeignKeys](fkColumns, func(x *modules.ForeignKeys, i int) bool {
		return x.PureName == table.PureName
	})

	grouped := lo.GroupBy[*modules.ForeignKeys, string](array, func(t *modules.ForeignKeys) string {
		return t.ConstraintName
	})

	var mapKeys []string
	for key, _ := range grouped {
		mapKeys = append(mapKeys, key)
	}

	return lo.Map(mapKeys, func(constraintName string, i int) map[string]interface{} {
		foreignKeys := map[string]interface{}{
			"constraintName": constraintName,
			"constraintType": "foreignKey",
		}
		if len(foreignKeys) > 0 {
			foreignKeys["constraintName"] = grouped[constraintName][0].ConstraintName
			//foreignKeys["schemaName"] = grouped[constraintName][0]
			foreignKeys["pureName"] = grouped[constraintName][0].PureName
			foreignKeys["refSchemaName"] = grouped[constraintName][0].RefTableName
			foreignKeys["updateAction"] = grouped[constraintName][0].UpdateAction
			foreignKeys["deleteAction"] = grouped[constraintName][0].DeleteAction
		}

		foreignKeys["columns"] = lo.Map[*modules.ForeignKeys, map[string]interface{}](grouped[constraintName], func(item *modules.ForeignKeys, i int) map[string]interface{} {
			return map[string]interface{}{
				"columnName":    item.ColumnName,
				"refColumnName": item.RefColumnName,
			}
		})

		return foreignKeys
	})
}

func CreateEmptyStructure() map[string]interface{} {
	return map[string]interface{}{
		"tables":      make([]map[string]interface{}, 0),
		"collections": make([]map[string]interface{}, 0),
		"views":       make([]map[string]interface{}, 0),
		"matviews":    make([]map[string]interface{}, 0),
		"functions":   make([]map[string]interface{}, 0),
		"procedures":  make([]map[string]interface{}, 0),
		"triggers":    make([]map[string]interface{}, 0),
		"schemas":     make([]map[string]interface{}, 0),
	}
}
