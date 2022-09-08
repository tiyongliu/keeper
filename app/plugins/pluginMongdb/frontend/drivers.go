package frontend

func Driver() []map[string]interface{} {
	return []map[string]interface{}{
		{
			"databaseEngineTypes": []string{"document"},
			"title":               "MongoDB",
			"engine":              "mongo",
			"defaultPort":         27017,
		},
	}
}
