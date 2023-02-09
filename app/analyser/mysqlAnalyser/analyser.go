package mysqlAnalyser

import (
	"fmt"
	"github.com/samber/lo"
	"keeper/app/analyser"
	sql2 "keeper/app/analyser/mysqlAnalyser/sql"
	"keeper/app/db"
	"keeper/app/db/adapter/mysql"
	"keeper/app/db/standard/modules"
	"keeper/app/pkg/logger"
	"keeper/app/utility"
	"reflect"
	"strconv"
	"strings"
)

var sql map[string]string

func init() {
	sql = map[string]string{
		"columns":                sql2.ColumnsSQL(),
		"tables":                 sql2.TablesSQL(),
		"primaryKeys":            sql2.PrimaryKeysSQL(),
		"foreignKeys":            sql2.ForeignKeysSQL(),
		"tableModifications":     sql2.TableModificationsSQL(),
		"views":                  sql2.ViewsSQL(),
		"programmables":          sql2.ProgrammablesSQL(),
		"procedureModifications": sql2.ProcedureModificationsSQL(),
		"functionModifications":  sql2.FunctionModificationsSQL(),
		"indexes":                sql2.IndexesSQL(),
		"uniqueNames":            sql2.UniqueNamesSQL(),
	}
}

type Analyser struct {
	Driver           db.Session
	DatabaseName     string
	DatabaseAnalyser *analyser.DatabaseAnalyser
}

func NewAnalyser(driver db.Session, database string) *Analyser {
	return &Analyser{
		Driver:           driver,
		DatabaseName:     database,
		DatabaseAnalyser: analyser.NewDatabaseAnalyser(driver),
	}
}

func (as *Analyser) CreateQuery(resFileName string, typeFields []string) string {
	res := sql[resFileName]
	res = strings.Replace(res, "#DATABASE#", as.DatabaseName, -1)

	return analyser.NewDatabaseAnalyser(as.Driver).CreateQuery(res, typeFields)
}

func (as *Analyser) RunAnalysis() map[string]interface{} {
	driver, ok := as.Driver.(*mysql.Source)
	if !ok && driver == nil {
		return nil
	}

	tables, _ := driver.Tables(as.CreateQuery("tables", []string{"tables"}))
	columns, _ := driver.Columns(as.CreateQuery("columns", []string{"tables", "views"}))

	pkColumns, err := driver.PrimaryKeys(as.CreateQuery("primaryKeys", []string{"tables"}))
	if err != nil {
		logger.Errorf("Error running analyser query %v", err)

		pkColumns = &modules.MysqlRowsResult{Rows: []*modules.PrimaryKey{}}
	}

	fkColumns, err := driver.ForeignKeys(as.CreateQuery("foreignKeys", []string{"tables"}))
	if err != nil {
		logger.Errorf("Error running analyser query %v", err)
		fkColumns = &modules.MysqlRowsResult{Rows: []*modules.ForeignKeys{}}
	}

	views, err := driver.Views(as.CreateQuery("views", []string{"views"}))
	if err != nil {
		logger.Errorf("Error running analyser query %v", err)
		views = &modules.MysqlRowsResult{Rows: []*modules.View{}}
	}

	viewTexts := getViewTexts(views.Rows.([]*modules.View))

	programmables, err := driver.Programmables(as.CreateQuery("programmables", []string{"procedures", "functions"}))
	if err != nil {
		logger.Errorf("Error running analyser query %v", err)
		programmables = &modules.MysqlRowsResult{Rows: []*modules.Programmable{}}
	}

	indexes, err := driver.Indexes(as.CreateQuery("indexes", []string{"tables"}))
	if err != nil {
		logger.Errorf("Error running analyser query %v", err)
		indexes = &modules.MysqlRowsResult{Rows: []*modules.Indexe{}}
	}

	uniqueNames, err := driver.UniqueNames(as.CreateQuery("uniqueNames", []string{"tables"}))
	if err != nil {
		uniqueNames = &modules.MysqlRowsResult{Rows: []*modules.UniqueName{}}
	}

	respAnalyser := make(map[string]interface{})
	if tables != nil {
		respAnalyser["tables"] = lo.Map(tables.Rows.([]*modules.Table), func(table *modules.Table, i int) map[string]interface{} {
			return map[string]interface{}{
				"pureName":      table.PureName,
				"tableRowCount": strconv.Itoa(table.TableRowCount),
				"modifyDate":    table.ModifyDate,
				"objectId":      table.PureName,
				"contentHash":   table.ModifyDate,
				"columns":       transformTablesColumns(table, columns.Rows.([]*modules.TableColumn)),
				"primaryKey":    analyser.ExtractPrimaryKeys(table, pkColumns.Rows.([]*modules.PrimaryKey)),
				"foreignKeys":   analyser.ExtractForeignKeys(table, fkColumns.Rows.([]*modules.ForeignKeys)),
				"indexes":       transformTablesIndexes(table, indexes.Rows.([]*modules.Indexe), uniqueNames.Rows.([]*modules.UniqueName)),
				"uniques":       transformTablesUniques(table, indexes.Rows.([]*modules.Indexe), uniqueNames.Rows.([]*modules.UniqueName)),
			}
		})
	}

	respAnalyser["views"] = lo.Map(views.Rows.([]*modules.View), func(view *modules.View, i int) map[string]interface{} {
		return map[string]interface{}{
			"pureName":       view.PureName,
			"modifyDate":     view.ModifyDate,
			"objectId":       view.PureName,
			"contentHash":    view.ModifyDate,
			"columns":        transformViewColumns(view, columns.Rows.([]*modules.TableColumn)),
			"createSql":      viewTexts, //todo 后期需要完善这里，目前没有数据无法编写
			"requiresFormat": true,
		}
	})

	respAnalyser["procedures"] = transformProcedures(programmables.Rows.([]*modules.Programmable))

	respAnalyser["functions"] = transformFunctions(programmables.Rows.([]*modules.Programmable))

	return respAnalyser
}

func transformTablesIndexes(table *modules.Table, indexesRows []*modules.Indexe, uniqueNamesRows []*modules.UniqueName) []map[string]interface{} {
	filters := lo.Filter[*modules.Indexe](indexesRows, func(idx *modules.Indexe, _ int) bool {
		existing, _ := lo.Find[*modules.UniqueName](uniqueNamesRows, func(x *modules.UniqueName) bool {
			return x.ConstraintName == idx.ConstraintName
		})

		return idx.TableName == table.PureName && existing == nil
	})

	uniqBy := lo.UniqBy[*modules.Indexe](filters, func(t *modules.Indexe) string {
		return t.ConstraintName
	})

	return lo.Map(uniqBy, func(idx *modules.Indexe, i int) map[string]interface{} {
		cols := lo.Filter[*modules.Indexe](indexesRows, func(col *modules.Indexe, _ int) bool {
			return col.TableName == idx.TableName && col.ConstraintName == idx.ConstraintName
		})
		return map[string]interface{}{
			"constraintName": idx.ConstraintName,
			"indexType":      idx.IndexType,
			"isUnique":       !idx.NonUnique,
			"columns": lo.Map(cols, func(col *modules.Indexe, _ int) map[string]interface{} {
				pick := make(map[string]interface{})
				pick["columnName"] = col.ColumnName
				return pick
			}),
		}
	})
}

func transformTablesColumns(table *modules.Table, columnsRows []*modules.TableColumn) []*modules.TransformColumnInfo {
	return getColumnInfo(lo.Filter[*modules.TableColumn](columnsRows, func(col *modules.TableColumn, _ int) bool {
		return col.PureName == table.PureName
	}))
}

func transformTablesUniques(table *modules.Table, indexesRows []*modules.Indexe, uniqueNamesRows []*modules.UniqueName) []map[string]interface{} {
	filters := lo.Filter[*modules.Indexe](indexesRows, func(idx *modules.Indexe, _ int) bool {
		existing, _ := lo.Find[*modules.UniqueName](uniqueNamesRows, func(x *modules.UniqueName) bool {
			return x.ConstraintName == idx.ConstraintName
		})

		return idx.TableName == table.PureName && existing != nil
	})

	uniqBy := lo.UniqBy[*modules.Indexe](filters, func(t *modules.Indexe) string {
		return t.ConstraintName
	})

	return lo.Map(uniqBy, func(idx *modules.Indexe, i int) map[string]interface{} {
		cols := lo.Filter[*modules.Indexe](indexesRows, func(col *modules.Indexe, _ int) bool {
			return col.TableName == idx.TableName && col.ConstraintName == idx.ConstraintName
		})

		return map[string]interface{}{
			"constraintName": idx.ConstraintName,
			"columns": lo.Map(cols, func(col *modules.Indexe, _ int) map[string]interface{} {
				pick := make(map[string]interface{})
				pick["columnName"] = col.ColumnName
				return pick
			}),
		}
	})

}

func transformViewColumns(view *modules.View, columnsRows []*modules.TableColumn) []*modules.TransformColumnInfo {
	return getColumnInfo(lo.Filter[*modules.TableColumn](columnsRows, func(col *modules.TableColumn, _ int) bool {
		return col.PureName == view.PureName
	}))
}

func transformProcedures(programmablesRows []*modules.Programmable) []map[string]interface{} {
	programmables := lo.Filter[*modules.Programmable](programmablesRows, func(x *modules.Programmable, _ int) bool {
		return x != nil && x.ObjectType == "PROCEDURE"
	})

	return lo.Map(programmables, func(x *modules.Programmable, i int) map[string]interface{} {
		return map[string]interface{}{
			"pureName":          x.PureName,
			"modifyDate":        x.ModifyDate,
			"returnDataType":    x.ReturnDataType,
			"routineDefinition": x.RoutineDefinition,
			"isDeterministic":   x.IsDeterministic,
			"createSql":         nil, //todo 后期需要完善这里，目前没有数据无法编写
			"objectId":          x.PureName,
			"contentHash":       x.ModifyDate,
		}
	})
}

func transformFunctions(programmablesRows []*modules.Programmable) []map[string]interface{} {
	functions := lo.Filter[*modules.Programmable](programmablesRows, func(x *modules.Programmable, _ int) bool {
		return x != nil && x.ObjectType == "FUNCTION"
	})

	return lo.Map(functions, func(x *modules.Programmable, i int) map[string]interface{} {
		return map[string]interface{}{
			"pureName":          x.PureName,
			"modifyDate":        x.ModifyDate,
			"returnDataType":    x.ReturnDataType,
			"routineDefinition": x.RoutineDefinition,
			"isDeterministic":   x.IsDeterministic,
			"createSql":         nil, //todo 后期需要完善这里，目前没有数据无法编写
			"objectId":          x.PureName,
			"contentHash":       x.ModifyDate,
		}
	})
}

func getColumnInfo(filter []*modules.TableColumn) []*modules.TransformColumnInfo {
	return lo.Map(filter, func(col *modules.TableColumn, i int) *modules.TransformColumnInfo {
		columnTypeTokens := lo.Map[string](strings.Split(col.ColumnType, " "), func(x string, i int) string {
			return strings.ToLower(strings.TrimSpace(x))
		})

		fullDataType := col.DataType
		if col.CharMaxLength != nil && utility.IsTypeString(col.DataType) {
			fullDataType = fmt.Sprintf("%s(%d)", col.DataType, *col.CharMaxLength)
		}

		if col.NumericPrecision != nil && col.NumericScale != nil && utility.IsTypeNumeric(col.DataType) {
			fullDataType = fmt.Sprintf("%s(%d,%d)", col.DataType, *col.NumericPrecision, *col.NumericScale)
		}

		var extra string
		utility.WithRecover(func() {
			if col.Extra != nil {
				extra = fmt.Sprintf("%s", col.Extra)
			}
		}, func(err error) {
			logger.Errorf("col.Extra expected string, got %s, conversion filed %v", reflect.ValueOf(col.Extra).Kind().String(), err)
		})

		return &modules.TransformColumnInfo{
			NotNull:       col.IsNullable == "" || strings.ToLower(col.IsNullable) == "no",
			AutoIncrement: !!(extra != "" && strings.Contains(strings.ToLower(extra), "auto_increment")),
			ColumnName:    col.ColumnName,
			ColumnComment: col.ColumnComment,
			DataType:      fullDataType,
			DefaultValue:  col.DefaultValue,
			IsUnsigned:    lo.Contains(columnTypeTokens, "unsigned"),
			IsZerofill:    lo.Contains(columnTypeTokens, "zerofill"),
		}
	})
}

func getViewTexts(viewsRows []*modules.View) interface{} {
	allViewNames := lo.Map[*modules.View, *modules.ViewTexts](viewsRows, func(x *modules.View, i int) *modules.ViewTexts {
		return &modules.ViewTexts{PureName: x.PureName}
	})
	getRequestedViewNames(allViewNames)
	return nil
}

func getRequestedViewNames(allViewNames []*modules.ViewTexts) {

}

func safeQuery[T any](result *modules.MysqlRowsResult, err error) *modules.MysqlRowsResult {
	if err != nil {
		logger.Errorf("Error running analyser query %v", err)
		result.Rows = make([]T, 0)
	}
	return result
}
