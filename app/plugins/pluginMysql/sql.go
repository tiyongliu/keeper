package pluginMysql

func columnsSQL() string {
	return `
		SELECT
            TABLE_NAME AS pureName,
			COLUMN_NAME AS columnName,
			IS_NULLABLE AS isNullable,
			DATA_TYPE AS dataType,
			CHARACTER_MAXIMUM_LENGTH AS charMaxLength,
			NUMERIC_PRECISION AS numericPrecision,
			NUMERIC_SCALE AS numericScale,
			COLUMN_DEFAULT AS defaultValue,
			COLUMN_COMMENT AS columnComment,
			COLUMN_TYPE AS columnType,
			EXTRA AS extra 
		FROM
			information_schema.COLUMNS 
		WHERE
			TABLE_SCHEMA = ? 
			AND TABLE_NAME = ? 
		ORDER BY
			ORDINAL_POSITION`
}

func foreignKeysSQL() string {
	return `select 
		REFERENTIAL_CONSTRAINTS.CONSTRAINT_NAME as constraintName,
		REFERENTIAL_CONSTRAINTS.TABLE_NAME as pureName,
		REFERENTIAL_CONSTRAINTS.UPDATE_RULE as updateAction,
		REFERENTIAL_CONSTRAINTS.DELETE_RULE as deleteAction,
		REFERENTIAL_CONSTRAINTS.REFERENCED_TABLE_NAME as refTableName,
		KEY_COLUMN_USAGE.COLUMN_NAME as columnName,
		KEY_COLUMN_USAGE.REFERENCED_COLUMN_NAME as refColumnName
	from INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS
	inner join INFORMATION_SCHEMA.KEY_COLUMN_USAGE 
		on REFERENTIAL_CONSTRAINTS.TABLE_NAME = KEY_COLUMN_USAGE.TABLE_NAME 
		and REFERENTIAL_CONSTRAINTS.CONSTRAINT_NAME = KEY_COLUMN_USAGE.CONSTRAINT_NAME
		and REFERENTIAL_CONSTRAINTS.CONSTRAINT_SCHEMA = KEY_COLUMN_USAGE.CONSTRAINT_SCHEMA
	where REFERENTIAL_CONSTRAINTS.CONSTRAINT_SCHEMA = ? and REFERENTIAL_CONSTRAINTS.TABLE_NAME =?
	order by KEY_COLUMN_USAGE.ORDINAL_POSITION`
}

func functionModificationSQL() string {
	return `SHOW FUNCTION STATUS WHERE Db = ?`
}

func indexesSQL() string {
	return `    SELECT 
		INDEX_NAME AS constraintName,
		TABLE_NAME AS tableName,
		COLUMN_NAME AS columnName,
		INDEX_TYPE AS indexType,
		NON_UNIQUE AS nonUnique
	FROM INFORMATION_SCHEMA.STATISTICS
	WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? AND INDEX_NAME != 'PRIMARY'
	ORDER BY SEQ_IN_INDEX`
}

func primaryKeysSQL() string {
	return `SELECT
		TABLE_CONSTRAINTS.CONSTRAINT_NAME AS constraintName,
		TABLE_CONSTRAINTS.TABLE_NAME AS pureName,
		KEY_COLUMN_USAGE.COLUMN_NAME AS columnName 
	FROM
		INFORMATION_SCHEMA.TABLE_CONSTRAINTS
		INNER JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE ON TABLE_CONSTRAINTS.TABLE_NAME = KEY_COLUMN_USAGE.TABLE_NAME 
		AND TABLE_CONSTRAINTS.CONSTRAINT_NAME = KEY_COLUMN_USAGE.CONSTRAINT_NAME 
		AND TABLE_CONSTRAINTS.CONSTRAINT_SCHEMA = KEY_COLUMN_USAGE.CONSTRAINT_SCHEMA 
	WHERE
		TABLE_CONSTRAINTS.CONSTRAINT_SCHEMA = ?
		AND TABLE_CONSTRAINTS.TABLE_NAME = ? 
		AND TABLE_CONSTRAINTS.CONSTRAINT_TYPE = 'PRIMARY KEY' 
	ORDER BY
		KEY_COLUMN_USAGE.ORDINAL_POSITION`
}

func procedureModificationsSQL() string {
	return `SHOW PROCEDURE STATUS WHERE Db = ?`
}

func programmablesSQL() string {
	return `select 
		ROUTINE_NAME as pureName,
		ROUTINE_TYPE as objectType,
		COALESCE(LAST_ALTERED, CREATED) as modifyDate,
		DATA_TYPE AS returnDataType,
		ROUTINE_DEFINITION as routineDefinition,
		IS_DETERMINISTIC as isDeterministic
	from information_schema.routines
	where ROUTINE_SCHEMA = ? and ROUTINE_NAME =? `
}

func tableModificationsSQL() string {
	return `select 
	TABLE_NAME as pureName, 
	TABLE_TYPE as objectType,
	TABLE_ROWS as tableRowCount,
	case when ENGINE='InnoDB' then CREATE_TIME else coalesce(UPDATE_TIME, CREATE_TIME) end as modifyDate 
from information_schema.tables 
where TABLE_SCHEMA = ? and TABLE_TYPE='BASE TABLE' and TABLE_NAME =?`
}

func tablesSQL() string {
	return `select
		TABLE_NAME as pureName,
		TABLE_ROWS as tableRowCount,
		case when ENGINE='InnoDB' then CREATE_TIME else coalesce(UPDATE_TIME, CREATE_TIME) end as modifyDate
	from information_schema.tables
	where TABLE_SCHEMA = ? and TABLE_TYPE='BASE TABLE' and TABLE_NAME =?`
}

func uniqueNamesSQL() string {
	return `select CONSTRAINT_NAME as constraintName
    from information_schema.TABLE_CONSTRAINTS
    where CONSTRAINT_SCHEMA = ? and constraint_type = 'UNIQUE'`
}

func viewsSQL() string {
	return `select 
		TABLE_NAME as pureName, 
		coalesce(UPDATE_TIME, CREATE_TIME) as modifyDate
	from information_schema.tables 
	where TABLE_SCHEMA = ? and TABLE_NAME = ? and TABLE_TYPE = 'VIEW'`
}
