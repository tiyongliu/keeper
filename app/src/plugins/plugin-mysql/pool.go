package plugin_mysql

import (
	"bytes"
	"dbbox/app/src/modules"
	"dbbox/app/src/pkg/standard"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

type SimpleSettingMysql struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     string `json:"port"`
}

func NewSimpleMysqlPool(setting *modules.SimpleSettingMysql) (standard.SqlStandard, error) {
	if setting == nil {
		return nil, fmt.Errorf("setting is nil")
	}
	if setting.Username == "" {
		return nil, fmt.Errorf("lack of setting.Username")
	}
	if setting.Password == "" {
		return nil, fmt.Errorf("lack of setting.Password")
	}
	if setting.Host == "" {
		return nil, fmt.Errorf("lack of setting.Host")
	}
	if setting.Port == "" {
		return nil, fmt.Errorf("lack of setting.Port")
	}

	var buf bytes.Buffer
	buf.WriteString(setting.Username)
	buf.WriteString(":")
	buf.WriteString(setting.Password)
	buf.WriteString("@tcp(")
	buf.WriteString(setting.Host)
	buf.WriteString(":")
	buf.WriteString(setting.Port)
	buf.WriteString(")/")
	buf.WriteString("mysql")
	buf.WriteString("?timeout=3s&readTimeout=5s")
	buf.WriteString("&loc=Local")
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,       // Disable color
		},
	)

	db, err := gorm.Open(mysql.New(mysql.Config{
		DriverName: "mysql",
		DSN:        buf.String(),
	}), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return nil, err
	}

	return &MysqlDrivers{db}, nil
}
