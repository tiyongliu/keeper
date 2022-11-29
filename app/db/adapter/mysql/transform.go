package mysql

import "keeper/app/db/standard/modules"

func transformMysqlDatabases(list []string) (lastDatabases []*modules.MysqlDatabase) {
	for _, value := range list {
		lastDatabases = append(lastDatabases, &modules.MysqlDatabase{Name: value})
	}

	return lastDatabases
}
