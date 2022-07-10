package plugin_mysql

import (
	"fmt"
	"gorm.io/gorm"
	"keeper/app/code"
	"keeper/app/modules"
	"keeper/app/pkg/logger"
	"keeper/app/pkg/standard"
	"regexp"
)

type MysqlDrivers struct {
	DB *gorm.DB
}

func NewMysql() standard.SqlStandard {
	return &MysqlDrivers{}
}

func (mysql *MysqlDrivers) Dialect() string {
	return code.MYSQLALIAS
}

func (mysql *MysqlDrivers) Connect() interface{} {
	return nil
}

func (mysql *MysqlDrivers) GetPoolInfo() interface{} {
	return mysql.DB
}

func (mysql *MysqlDrivers) GetVersion() (interface{}, error) {
	var rows []string
	err := mysql.DB.Raw("select version()").Scan(&rows).Error
	if err != nil {
		logger.Errorf("get mysql version failed: %v", err)
		return nil, err
	}

	if len(rows) > 0 && rows[0] != "" {
		subMath := regexp.MustCompile("(.*)-MariaDB-").FindAllSubmatch([]byte(rows[0]), -1)
		if len(subMath) >= 1 {
			return &standard.VersionMsg{
				Version:     rows[0],
				VersionText: fmt.Sprintf("MariaDB %s", subMath[0]),
			}, nil
		}
	}

	return &standard.VersionMsg{
		Version:     rows[0],
		VersionText: fmt.Sprintf("MySQL %s", rows[0]),
	}, nil
}

func (mysql *MysqlDrivers) ListDatabases() (interface{}, error) {
	var rows []string
	err := mysql.DB.Raw("SHOW DATABASES").Scan(&rows).Error
	if err != nil {
		logger.Errorf("get mysql lastDatabases failed: %v", err)
		return nil, err
	}
	return TransformListDatabases(rows), nil
}

func (mysql *MysqlDrivers) Close() error {
	db, err := mysql.DB.DB()
	if err != nil {
		return err
	}
	return db.Close()
}

func (mysql *MysqlDrivers) Tables() (interface{}, error) {
	var tableSchemas []*modules.MysqlTableSchema
	rows, err := mysql.DB.Raw(tableModifications(), "shop_go").Rows()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var pureName, objectType, modifyDate string
		var tableRowCount int
		if err := rows.Scan(&pureName, &objectType, &tableRowCount, &modifyDate); err != nil {
			return nil, err
		}

		tableSchemas = append(tableSchemas, &modules.MysqlTableSchema{
			PureName:      pureName,
			ObjectType:    objectType,
			TableRowCount: tableRowCount,
			ModifyDate:    modifyDate,
		})
	}
	return tableSchemas, nil
}

func (mysql *MysqlDrivers) Columns(databaseName, tableName string) (interface{}, error) {
	var columns []*modules.MysqlTableColumn
	rows, err := mysql.DB.Raw(informationSchemaColumns(), databaseName, tableName).Rows()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var pureName, columnName, isNullable, dataType, columnType, extra string
		var charMaxLength, numericPrecision, numericScale *int
		var defaultValue, columnComment interface{}
		if err := rows.Scan(&pureName,
			&columnName,
			&isNullable,
			&dataType,
			&charMaxLength,
			&numericPrecision,
			&numericScale,
			&defaultValue,
			&columnComment,
			&columnType,
			&extra,
		); err != nil {
			return nil, err
		}

		columns = append(columns, &modules.MysqlTableColumn{
			PureName:         pureName,
			ColumnName:       columnName,
			IsNullable:       isNullable,
			DAtaType:         dataType,
			CharMaxLength:    charMaxLength,
			NumericPrecision: numericPrecision,
			NumericScale:     numericScale,
			DefaultValue:     defaultValue,
			ColumnComment:    columnComment,
			ColumnType:       columnType,
			EXTRA:            extra,
		})
	}
	return columns, nil
}

func (mysql *MysqlDrivers) PrimaryKeys(databaseName, tableName string) (interface{}, error) {
	var primaryKeys []*modules.MysqlPrimaryKey
	rows, err := mysql.DB.Raw(informationSchemaColumnsPrimaryKeys(), databaseName, tableName).Rows()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var constraintName, pureName, columnName string
		if err := rows.Scan(&constraintName, &pureName, &columnName); err != nil {
			return nil, err
		}

		primaryKeys = append(primaryKeys, &modules.MysqlPrimaryKey{
			ConstraintName: constraintName,
			PureName:       pureName,
			ColumnName:     columnName,
		})
	}

	return primaryKeys, nil
}

/*
[2022-07-10 20:53:57.104][localhost_3306][000015][MYSQL]
SHOW TABLE STATUS LIKE 'tz_user'
Time: 0.001s

[2022-07-10 20:53:57.106][localhost_3306][000015][MYSQL]
SHOW CREATE TABLE `tz_user`
Time: 0.000s

[2022-07-10 20:53:57.107][localhost_3306][000015][MYSQL]
SHOW FULL COLUMNS FROM `tz_user`
Time: 0.002s

[2022-07-10 20:53:57.112][localhost_3306][000015][MYSQL]
SHOW INDEX FROM `tz_user`
Time: 0.002s

[2022-07-10 20:53:57.115][localhost_3306][000015][MYSQL]
SELECT action_order, event_object_table, trigger_name, event_manipulation, event_object_table, definer, action_statement, action_timing FROM information_schema.triggers WHERE BINARY event_object_schema = 'shop_go' AND BINARY event_object_table = 'tz_user' ORDER BY event_object_table
Time: 0.000s

[2022-07-10 20:53:57.116][localhost_3306][000015][MYSQL]
SELECT TABLE_NAME, PARTITION_NAME, SUBPARTITION_NAME, PARTITION_METHOD, SUBPARTITION_METHOD, PARTITION_EXPRESSION, SUBPARTITION_EXPRESSION, PARTITION_DESCRIPTION, PARTITION_COMMENT, NODEGROUP, TABLESPACE_NAME FROM information_schema.PARTITIONS WHERE NOT ISNULL(PARTITION_NAME) AND TABLE_SCHEMA LIKE 'shop_go' AND TABLE_NAME LIKE 'tz_user' ORDER BY TABLE_NAME, PARTITION_NAME, PARTITION_ORDINAL_POSITION, SUBPARTITION_ORDINAL_POSITION
Time: 0.001s

[2022-07-10 20:53:57.171][localhost_3306][000015][MYSQL]
SHOW DATABASES
Time: 0.000s
*/
