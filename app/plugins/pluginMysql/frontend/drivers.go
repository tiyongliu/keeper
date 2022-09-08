package frontend

import "keeper/app/utility"

func mysqlDriverBase() map[string]interface{} {
	return map[string]interface{}{
		"databaseEngineTypes": []string{"sql"},
		"defaultPort":         3306,
		"readOnlySessions":    true,
	}
}

func mysqlDriver() map[string]interface{} {
	return utility.MergeUnknownMaps(mysqlDriverBase(), map[string]interface{}{
		"title":  "MySQL",
		"engine": "mysql",
	})
}

func mariaDriver() map[string]interface{} {
	return utility.MergeUnknownMaps(mysqlDriverBase(), map[string]interface{}{
		"title":  "MariaDB",
		"engine": "mariadb",
	})
}

func Driver() []map[string]interface{} {
	return []map[string]interface{}{mysqlDriver(), mariaDriver()}
}
