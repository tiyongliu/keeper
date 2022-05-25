package plugin_mysql

import (
	"dbbox/app/src/pkg/logger"
	"dbbox/app/src/pkg/standard"
	"fmt"
	"regexp"
)

type MysqlDrivers struct {
}

func NewMysql() standard.SqlStandard {
	return &MysqlDrivers{}
}

func (mysql *MysqlDrivers) GetVersion() (interface{}, error) {
	connect, err := NewConnection("127.0.0.1", "3306", "root", "123456", "mysql")
	if err != nil {
		return nil, err
	}

	var rows []string

	err = connect.Raw("select version()").Scan(&rows).Error
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
