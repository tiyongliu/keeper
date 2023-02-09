package sql

func UniqueNamesSQL() string {
	return `select CONSTRAINT_NAME as constraintName
    from information_schema.TABLE_CONSTRAINTS
    where CONSTRAINT_SCHEMA = '#DATABASE#' and constraint_type = 'UNIQUE'`
}
