package plugin_mysql

import (
	"fmt"
	"gorm.io/gorm"
	"keeper/app/code"
	"keeper/app/pkg/logger"
	"keeper/app/pkg/standard"
	"keeper/app/tools"
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

type TableSchema struct {
	PureName      string `json:"pureName"`
	ObjectType    string `json:"objectType"`
	TableRowCount int    `json:"tableRowCount"`
	ModifyDate    string `json:"modifyDate"`
}

func (mysql *MysqlDrivers) Tables() (interface{}, error) {
	var tableSchemas []*TableSchema
	rows, err := mysql.DB.Raw(tableModifications(), "shop_go").Rows()
	if err != nil {
		logger.Errorf("err: %v", err)
		return nil, err
	}

	for rows.Next() {
		var pureName, objectType, modifyDate string
		var tableRowCount int
		if err := rows.Scan(&pureName, &objectType, &tableRowCount, &modifyDate); err != nil {
			logger.Errorf("err: %v", err)
			return nil, err
		}

		tableSchemas = append(tableSchemas, &TableSchema{
			PureName:      pureName,
			ObjectType:    objectType,
			TableRowCount: tableRowCount,
			ModifyDate:    modifyDate,
		})
	}
	logger.Infof("%s", tools.ToJsonStr(tableSchemas))
	return tableSchemas, nil
}
