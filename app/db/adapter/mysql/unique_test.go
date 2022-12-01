package mysql

import (
	"fmt"
	"keeper/app/db"
	"keeper/app/utility"
	"testing"
)

func TestUniqueNames(t *testing.T) {
	getDevice(func(session db.Session) {
		driver, ok := session.(*Source)
		if !ok && driver == nil {
			return
		}

		names, err := driver.UniqueNames(`select CONSTRAINT_NAME as constraintName
			from information_schema.TABLE_CONSTRAINTS
		    where CONSTRAINT_SCHEMA = 'yami_shops' and constraint_type = 'UNIQUE'`)

		fmt.Printf("RES: %s, err: [%v]", utility.ToJsonStr(names.Rows), err)
	})
}

func TestIndexes(t *testing.T) {
	getDevice(func(session db.Session) {
		driver, ok := session.(*Source)
		if !ok && driver == nil {
			return
		}

		indexes, err := driver.Indexes(`SELECT
	       INDEX_NAME AS constraintName,
	       TABLE_NAME AS tableName,
	       COLUMN_NAME AS columnName,
	       INDEX_TYPE AS indexType,
	       NON_UNIQUE AS nonUnique
	   FROM INFORMATION_SCHEMA.STATISTICS
	   WHERE TABLE_SCHEMA = 'crmeb' AND TABLE_NAME  is not null AND INDEX_NAME != 'PRIMARY'
	   ORDER BY SEQ_IN_INDEX`)
		fmt.Printf("RES: %s, err: [%v]", utility.ToJsonStr(indexes), err)
	})
}

func TestTables(t *testing.T) {
	getDevice(func(session db.Session) {
		driver, ok := session.(*Source)
		if !ok && driver == nil {
			return
		}
		tables_, err := driver.Tables(`select
		        TABLE_NAME as pureName,
		        TABLE_ROWS as tableRowCount,
		        case when ENGINE='InnoDB' then CREATE_TIME else coalesce(UPDATE_TIME, CREATE_TIME) end as modifyDate
		from information_schema.tables
		where TABLE_SCHEMA = 'crmeb' and TABLE_TYPE='BASE TABLE' and TABLE_NAME  is not null`)
		fmt.Printf("RES: %s, err: [%v]", utility.ToJsonStr(tables_), err)
	})
}

func TestColumns(t *testing.T) {
	getDevice(func(session db.Session) {
		driver, ok := session.(*Source)
		if !ok && driver == nil {
			return
		}
		columns_, err := driver.Columns(`select
        TABLE_NAME as pureName,
        COLUMN_NAME as columnName,
        IS_NULLABLE as isNullable,
        DATA_TYPE as dataType,
        CHARACTER_MAXIMUM_LENGTH as charMaxLength,
        NUMERIC_PRECISION as numericPrecision,
        NUMERIC_SCALE as numericScale,
        COLUMN_DEFAULT as defaultValue,
        COLUMN_COMMENT as columnComment,
        COLUMN_TYPE as columnType,
        EXTRA as extra
from INFORMATION_SCHEMA.COLUMNS
where TABLE_SCHEMA = 'ssodb' and TABLE_NAME  is not null
order by ORDINAL_POSITION`)
		fmt.Printf("RES: %s, err: [%v]", utility.ToJsonStr(columns_), err)
	})
}

func TestPrimaryKeys(t *testing.T) {
	getDevice(func(session db.Session) {
		driver, ok := session.(*Source)
		if !ok && driver == nil {
			return
		}

		keys_, err := driver.PrimaryKeys(`select
        TABLE_CONSTRAINTS.CONSTRAINT_NAME as constraintName,
        TABLE_CONSTRAINTS.TABLE_NAME as pureName,
        KEY_COLUMN_USAGE.COLUMN_NAME as columnName
from INFORMATION_SCHEMA.TABLE_CONSTRAINTS
inner join INFORMATION_SCHEMA.KEY_COLUMN_USAGE
        on TABLE_CONSTRAINTS.TABLE_NAME = KEY_COLUMN_USAGE.TABLE_NAME
                and TABLE_CONSTRAINTS.CONSTRAINT_NAME = KEY_COLUMN_USAGE.CONSTRAINT_NAME
                and TABLE_CONSTRAINTS.CONSTRAINT_SCHEMA = KEY_COLUMN_USAGE.CONSTRAINT_SCHEMA
where TABLE_CONSTRAINTS.CONSTRAINT_SCHEMA = 'ssodb' and TABLE_CONSTRAINTS.TABLE_NAME  is not null AND TABLE_CONSTRAINTS.CONSTRAINT_TYPE = 'PRIMARY KEY'
order by KEY_COLUMN_USAGE.ORDINAL_POSITION`)
		fmt.Printf("RES: %s, err: [%v]", utility.ToJsonStr(keys_), err)
	})
}

func TestForeignKeys(t *testing.T) {
	getDevice(func(session db.Session) {
		driver, ok := session.(*Source)
		if !ok && driver == nil {
			return
		}

		keys_, err := driver.ForeignKeys(`select
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
where REFERENTIAL_CONSTRAINTS.CONSTRAINT_SCHEMA = 'yami_shops' and REFERENTIAL_CONSTRAINTS.TABLE_NAME  is not null
order by KEY_COLUMN_USAGE.ORDINAL_POSITION`)
		fmt.Printf("RES: %s, err: [%v]", utility.ToJsonStr(keys_.Columns), err)
	})
}

func TestViews(t *testing.T) {
	getDevice(func(session db.Session) {
		driver, ok := session.(*Source)
		if !ok && driver == nil {
			return
		}

		views, err := driver.Views(`select
		        TABLE_NAME as pureName,
		    coalesce(UPDATE_TIME, CREATE_TIME) as modifyDate
		from information_schema.tables
		where TABLE_SCHEMA = 'crmeb' and TABLE_NAME  is not null and TABLE_TYPE = 'VIEW'`)
		fmt.Printf("RES: %s, err: [%v]", utility.ToJsonStr(views), err)
	})
}

func TestProgrammables(t *testing.T) {
	getDevice(func(session db.Session) {
		driver, ok := session.(*Source)
		if !ok && driver == nil {
			return
		}

		programmables, err := driver.Programmables(`select
    ROUTINE_NAME as pureName,
    ROUTINE_TYPE as objectType,
    COALESCE(LAST_ALTERED, CREATED) as modifyDate,
    DATA_TYPE AS returnDataType,
    ROUTINE_DEFINITION as routineDefinition,
    IS_DETERMINISTIC as isDeterministic
from information_schema.routines
where ROUTINE_SCHEMA = 'sys_region' and ROUTINE_NAME  is not null`)
		fmt.Printf("RES: %s, err: [%v]", utility.ToJsonStr(programmables), err)
	})
}
