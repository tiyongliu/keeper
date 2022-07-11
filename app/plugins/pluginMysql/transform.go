package pluginMysql

import "keeper/app/modules"

type DatabasesItem struct {
	Name string `json:"name"`
}

func TransformListDatabases(list []string) (lastDatabases []*modules.MysqlDatabase) {
	for _, value := range list {
		lastDatabases = append(lastDatabases, &modules.MysqlDatabase{Name: value})
	}

	return lastDatabases
}
