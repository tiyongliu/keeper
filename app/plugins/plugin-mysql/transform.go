package plugin_mysql

type DatabasesItem struct {
	Name string `json:"name"`
}

func TransformListDatabases(list []string) (lastDatabases []*DatabasesItem) {
	for _, value := range list {
		lastDatabases = append(lastDatabases, &DatabasesItem{Name: value})
	}

	return lastDatabases
}
