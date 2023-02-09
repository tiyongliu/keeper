package sql

func IndexesSQL() string {
	return `SELECT 
        INDEX_NAME AS constraintName,
        TABLE_NAME AS tableName,
        COLUMN_NAME AS columnName,
        INDEX_TYPE AS indexType,
        NON_UNIQUE AS nonUnique
    FROM INFORMATION_SCHEMA.STATISTICS
    WHERE TABLE_SCHEMA = '#DATABASE#' AND TABLE_NAME =OBJECT_ID_CONDITION AND INDEX_NAME != 'PRIMARY'
    ORDER BY SEQ_IN_INDEX`
}
