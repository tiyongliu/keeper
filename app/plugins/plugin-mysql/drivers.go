package plugin_mysql

import (
	"fmt"
	"gorm.io/gorm"
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
	return nil, nil
}

func (mysql *MysqlDrivers) Close() error {
	db, err := mysql.DB.DB()
	if err != nil {
		return err
	}
	return db.Close()
}
