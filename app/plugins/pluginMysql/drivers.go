package pluginMysql

import (
	"database/sql"
	"fmt"
	"gorm.io/gorm"
	"keeper/app/pkg/logger"
	"keeper/app/pkg/standard"
	"keeper/app/plugins/modules"
	"regexp"
)

type MysqlDrivers struct {
	DB *gorm.DB
}

func NewMysql() standard.SqlStandard {
	return &MysqlDrivers{}
}

func (mysql *MysqlDrivers) Dialect() string {
	return standard.MYSQLALIAS
}

func (mysql *MysqlDrivers) Connect() interface{} {
	return nil
}

func (mysql *MysqlDrivers) GetPoolInfo() interface{} {
	return mysql.DB
}

func (mysql *MysqlDrivers) GetVersion() (*standard.VersionMsg, error) {
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

func (mysql *MysqlDrivers) Ping() error {
	db, err := mysql.DB.DB()
	if err != nil {
		return err
	}
	return db.Ping()
}

func (mysql *MysqlDrivers) UniqueNames(sql string) (*modules.MysqlRowsResult, error) {
	rows, err := mysql.DB.Raw(sql).Rows()
	if err != nil {
		return nil, err
	}

	columns := getSqlColumns(rows)

	var name string
	var uniqueNames []*modules.UniqueName

	for rows.Next() {
		if err = rows.Scan(&name); err != nil {
			return nil, err
		} else {
			uniqueNames = append(uniqueNames, &modules.UniqueName{ConstraintName: name})
		}
	}

	return &modules.MysqlRowsResult{
		Rows:    uniqueNames,
		Columns: columns,
	}, nil
}

func (mysql *MysqlDrivers) Indexes(sql string) (*modules.MysqlRowsResult, error) {
	rows, err := mysql.DB.Raw(sql).Rows()
	if err != nil {
		return nil, err
	}

	columns := getSqlColumns(rows)

	var constraintName, tableName, columnName, indexType string
	var nonUnique bool
	var indexes []*modules.Indexe
	for rows.Next() {
		if err = rows.Scan(&constraintName, &tableName, &columnName, &indexType, &nonUnique); err != nil {
			return nil, err
		} else {
			indexes = append(indexes, &modules.Indexe{
				ConstraintName: constraintName,
				TableName:      tableName,
				ColumnName:     columnName,
				IndexType:      indexType,
				NonUnique:      nonUnique,
			})
		}
	}

	return &modules.MysqlRowsResult{
		Rows:    indexes,
		Columns: columns,
	}, nil
}

func (mysql *MysqlDrivers) Tables(sql string) (*modules.MysqlRowsResult, error) {
	rows, err := mysql.DB.Raw(sql).Rows()

	if err != nil {
		return nil, err
	}

	columns := getSqlColumns(rows)

	var pureName string
	var tableRowCount int
	var modifyDate string

	var tables []*modules.Table
	for rows.Next() {
		if err = rows.Scan(&pureName, &tableRowCount, &modifyDate); err != nil {
			return nil, err
		} else {
			tables = append(tables, &modules.Table{
				PureName:      pureName,
				TableRowCount: tableRowCount,
				ModifyDate:    modifyDate,
			})
		}
	}

	return &modules.MysqlRowsResult{
		Rows:    tables,
		Columns: columns,
	}, nil
}

func (mysql *MysqlDrivers) Columns(sql string) (*modules.MysqlRowsResult, error) {
	sqlQuery, err := getSqlQuery(mysql.DB, sql)
	if err != nil {
		return nil, err
	}

	var pureName, columnName, isNullable, dataType, columnComment, columnType string
	var charMaxLength, numericPrecision, numericScale *int
	var defaultValue, extra interface{}
	var tableColumns []*modules.TableColumn
	for sqlQuery.Rows.Next() {
		if err = sqlQuery.Rows.Scan(&pureName, &columnName, &isNullable, &dataType, &charMaxLength, &numericPrecision, &numericScale, &defaultValue, &columnComment, &columnType, &extra); err != nil {
			return nil, err
		} else {
			tableColumns = append(tableColumns, &modules.TableColumn{
				PureName:         pureName,
				ColumnName:       columnName,
				IsNullable:       isNullable,
				DataType:         dataType,
				CharMaxLength:    charMaxLength,
				NumericPrecision: numericPrecision,
				NumericScale:     numericScale,
				DefaultValue:     defaultValue,
				ColumnComment:    columnComment,
				ColumnType:       columnType,
				Extra:            extra,
			})
		}
	}

	return &modules.MysqlRowsResult{
		Rows:    tableColumns,
		Columns: sqlQuery.Columns,
	}, nil
}

type Query struct {
	Rows    *sql.Rows
	Columns []*modules.Column `json:"columns"`
}

func getSqlQuery(db *gorm.DB, sql string) (*Query, error) {
	rows, err := db.Raw(sql).Rows()
	if err != nil {
		return nil, nil
	}
	return &Query{Rows: rows, Columns: getSqlColumns(rows)}, nil
}

func getSqlColumns(rows *sql.Rows) (columns []*modules.Column) {
	rowsColumns, err := rows.Columns()
	if err != nil {
		return nil
	}
	for _, s := range rowsColumns {
		columns = append(columns, &modules.Column{ColumnName: s})
	}
	return columns
}

func (mysql *MysqlDrivers) PrimaryKeys(sql string) (*modules.MysqlRowsResult, error) {
	rows, err := mysql.DB.Raw(sql).Rows()
	if err != nil {
		return nil, err
	}

	columns := getSqlColumns(rows)

	var primaryKeys []*modules.PrimaryKey
	var constraintName, pureName, columnName string

	for rows.Next() {
		if err = rows.Scan(&constraintName, &pureName, &columnName); err != nil {
			return nil, err
		} else {
			primaryKeys = append(primaryKeys, &modules.PrimaryKey{
				ConstraintName: constraintName,
				PureName:       pureName,
				ColumnName:     columnName,
			})
		}
	}

	return &modules.MysqlRowsResult{Rows: primaryKeys, Columns: columns}, nil
}

func (mysql *MysqlDrivers) ForeignKeys(sql string) (*modules.MysqlRowsResult, error) {
	sqlQuery, err := getSqlQuery(mysql.DB, sql)
	if err != nil {
		return nil, err
	}

	var foreignKeys []*modules.ForeignKeys
	var constraintName, pureName, updateAction, deleteAction, refTableName, columnName, refColumnName string

	for sqlQuery.Rows.Next() {
		if err = sqlQuery.Rows.Scan(&constraintName,
			&pureName,
			&updateAction,
			&deleteAction,
			&refTableName,
			&columnName,
			&refColumnName); err != nil {
			return nil, err
		} else {
			foreignKeys = append(foreignKeys, &modules.ForeignKeys{
				ConstraintName: constraintName,
				PureName:       pureName,
				UpdateAction:   updateAction,
				DeleteAction:   deleteAction,
				RefTableName:   refTableName,
				ColumnName:     columnName,
				RefColumnName:  refColumnName,
			})
		}
	}

	return &modules.MysqlRowsResult{Rows: foreignKeys, Columns: sqlQuery.Columns}, nil
}

func (mysql *MysqlDrivers) Views(sql string) (*modules.MysqlRowsResult, error) {
	sqlQuery, err := getSqlQuery(mysql.DB, sql)
	if err != nil {
		return nil, err
	}

	var pureName, modifyDate string
	var views []*modules.View

	for sqlQuery.Rows.Next() {
		if err = sqlQuery.Rows.Scan(&pureName, &modifyDate); err != nil {
			return nil, err
		} else {
			views = append(views, &modules.View{PureName: pureName, ModifyDate: modifyDate})
		}
	}

	return &modules.MysqlRowsResult{
		Rows:    views,
		Columns: sqlQuery.Columns,
	}, nil
}

func (mysql *MysqlDrivers) Programmables(sql string) (*modules.MysqlRowsResult, error) {
	sqlQuery, err := getSqlQuery(mysql.DB, sql)
	if err != nil {
		return nil, err
	}

	var programmables []*modules.Programmable
	var pureName, objectType, modifyDate, returnDataType string
	var routineDefinition interface{}
	var isDeterministic bool
	for sqlQuery.Rows.Next() {
		if err = sqlQuery.Rows.Scan(&pureName, &modifyDate); err != nil {
			return nil, err
		} else {
			programmables = append(programmables, &modules.Programmable{
				PureName:          pureName,
				ObjectType:        objectType,
				ModifyDate:        modifyDate,
				ReturnDataType:    returnDataType,
				RoutineDefinition: routineDefinition,
				IsDeterministic:   isDeterministic,
			})
		}
	}

	return &modules.MysqlRowsResult{Rows: programmables, Columns: sqlQuery.Columns}, nil
}

func (mysql *MysqlDrivers) Query(sql string) (interface{}, error) {
	rows, err := mysql.DB.Raw(sql).Rows()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
	}
	return nil, nil
}
