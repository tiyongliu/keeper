package plugin_mysql

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewConnection(host, port, username, password, database string) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.New(mysql.Config{
		DriverName: "mysql",
		//ServerVersion:             "",
		DSN:                       fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?timeout=3s&readTimeout=5s", username, password, host, port, database),
		Conn:                      nil,
		SkipInitializeWithVersion: false,
		DefaultStringSize:         256,
		//DefaultDatetimePrecision:  nil,
		DisableDatetimePrecision: true,
		//DontSupportRenameIndex:    false,
		//DontSupportRenameColumn:   false,
		//DontSupportForShareClause: false,
	}))
	if err != nil {
		return nil, err
	}
	return db, nil
}
