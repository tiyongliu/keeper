package sql

func FunctionModificationsSQL() string {
	return `SHOW FUNCTION STATUS WHERE Db = '#DATABASE#'`
}
