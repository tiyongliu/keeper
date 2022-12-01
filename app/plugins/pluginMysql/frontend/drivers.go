package frontend

import (
	"github.com/samber/lo"
)

type ConnectionFieldValues struct {
	AuthType string `json:"authType"`
}

type PreConnectDriver interface {
}

func mysqlDriverBase() map[string]interface{} {
	return map[string]interface{}{
		"databaseEngineTypes":  []string{"sql"},
		"defaultPort":          3306,
		"readOnlySessions":     true,
		"supportsDatabaseDump": true,
		"authTypeLabel":        "Connection mode",
		"defaultAuthTypeName":  "hostPort",
	}
}

func valuesAuthType(field string, values *ConnectionFieldValues) bool {
	return values != nil && (values.AuthType == "socket" && lo.Contains([]string{"socketPath"}, field) ||
		values.AuthType != "socket" && lo.Contains([]string{"server", "port"}, field))
}

func mysqlDriver() map[string]interface{} {
	return lo.Assign(mysqlDriverBase(), dialect, map[string]interface{}{
		"title":  "MySQL",
		"engine": "mysql",
	})
}

func mariaDriver() map[string]interface{} {
	return lo.Assign(
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
	"rangeSelect":                true,
	"stringEscapeChar":           "\\",
	"fallbackDataType":           "longtext",
	"enableConstraintsPerTable":  false,
	"anonymousPrimaryKey":        true,
	"explicitDropConstraint":     true,
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
}
