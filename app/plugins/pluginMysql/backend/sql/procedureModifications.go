package sql

func ProcedureModificationsSQL() string {
	return `SHOW PROCEDURE STATUS WHERE Db = '#DATABASE#'`
}
