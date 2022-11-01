package frontend

import (
	"github.com/samber/lo"
	"keeper/app/utility"
	"strings"
)

type ConnectionFieldValues struct {
	AuthType string `json:"authType"`
}

func mysqlDriverBase() map[string]interface{} {
	return map[string]interface{}{
		"databaseEngineTypes": []string{"sql"},
		"showConnectionField": func(field string, values *ConnectionFieldValues) bool {
			return lo.Contains([]string{"authType", "user", "password", "defaultDatabase", "singleDatabase", "isReadOnly"}, field) ||
				valuesAuthType(field, values)
		},
		"defaultPort":          3306,
		"readOnlySessions":     true,
		"supportsDatabaseDump": true,
		"authTypeLabel":        "Connection mode",
		"defaultAuthTypeName":  "hostPort",
		"getNewObjectTemplates": func() []map[string]string {
			return []map[string]string{
				{"label": "New view", "sql": "CREATE VIEW myview\nAS\nSELECT * FROM table1"},
				{
					"label": "New procedure",
					"sql":   "DELIMITER //\n\nCREATE PROCEDURE myproc (IN arg1 INT)\nBEGIN\n  SELECT * FROM table1;\nEND\n\nDELIMITER ;",
				},
				{
					"label": "New function",
					"sql":   "CREATE FUNCTION myfunc (arg1 INT)\nRETURNS INT DETERMINISTIC\nRETURN 1",
				},
			}
		},
	}
}

func valuesAuthType(field string, values *ConnectionFieldValues) bool {
	return values != nil && (values.AuthType == "socket" && lo.Contains([]string{"socketPath"}, field) ||
		values.AuthType != "socket" && lo.Contains([]string{"server", "port"}, field))
}

func mysqlDriver() map[string]interface{} {
	return utility.MergeUnknownMaps(mysqlDriverBase(), dialect, map[string]interface{}{
		"title":  "MySQL",
		"engine": "mysql",
	})
}

func mariaDriver() map[string]interface{} {
	return utility.MergeUnknownMaps(
		mysqlDriverBase(),
		dialect,
		map[string]interface{}{
			"title":  "MariaDB",
			"engine": "mariadb",
		},
	)
}

func Driver() []map[string]interface{} {
	return []map[string]interface{}{mysqlDriver(), mariaDriver()}
}

var spatialTypes = []string{
	"POINT",
	"LINESTRING",
	"POLYGON",
	"GEOMETRY",
	"MULTIPOINT",
	"MULTILINESTRING",
	"MULTIPOLYGON",
	"GEOMCOLLECTION",
	"GEOMETRYCOLLECTION",
}

var dialect = map[string]interface{}{
	"rangeSelect":               true,
	"stringEscapeChar":          "\\",
	"fallbackDataType":          "longtext",
	"enableConstraintsPerTable": false,
	"anonymousPrimaryKey":       true,
	"explicitDropConstraint":    true,
	"quoteIdentifier":           func(s string) string { return "`" + s + "`" },

	"createColumn":               true,
	"dropColumn":                 true,
	"changeColumn":               true,
	"createIndex":                true,
	"dropIndex":                  true,
	"createForeignKey":           true,
	"dropForeignKey":             true,
	"createPrimaryKey":           true,
	"dropPrimaryKey":             true,
	"dropIndexContainsTableSpec": true,
	"createUnique":               true,
	"dropUnique":                 true,
	"createCheck":                true,
	"dropCheck":                  true,

	"dropReferencesWhenDropTable": false,

	"columnProperties": map[string]interface{}{
		"columnComment": true,
		"isUnsigned":    true,
		"isZerofill":    true,
	},

	"predefinedDataTypes": []string{
		"char(20)",
		"varchar(250)",
		"binary(250)",
		"varbinary(250)",
		"tinyblob",
		"tinytext",
		"text(1000)",
		"blob(1000)",
		"mediumtext",
		"mediumblob",
		"longtext",
		"longblob",
		"enum(val1,val2,val3)",
		"set(val1,val2,val3)",
		"bit(32)",
		"tinyint",
		"bool",
		"smallint",
		"mediumint",
		"int",
		"bigint",
		"float",
		"double",
		"decimal",
		"date",
		"datetime",
		"timestamp",
		"time",
		"year",
	},

	"createColumnViewExpression": func(columnName, dataType, source string, alias interface{}) map[string]interface{} {
		if dataType != "" && lo.Contains(spatialTypes, strings.ToUpper(dataType)) {
			m := map[string]interface{}{
				"exprType": "call",
				"func":     "ST_AsText",
				"args": []map[string]interface{}{
					{"exprType": "column", "columnName": columnName, "source": source},
				},
			}
			if alias != nil {
				m["alias"] = alias
			} else {
				m["alias"] = columnName
			}
			return m
		}
		return nil
	},
}
