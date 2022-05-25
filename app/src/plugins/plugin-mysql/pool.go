package plugin_mysql

import (
	"bytes"
	"dbbox/app/src/modules"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"net/url"
	"os"
	"strconv"
	"time"
)

type SimpleSettingMysql struct {
	Host            string `yaml:"host" json:"host"`
	Username        string `yaml:"username" json:"username"`
	Password        string `yaml:"password" json:"password"`
	DBName          string `yaml:"dBName" json:"dbName"`
	Charset         string `yaml:"charset" json:"charset"`
	MaxIdle         int    `yaml:"maxIdle" json:"maxIdle"`
	MaxOpen         int    `yaml:"maxOpen" json:"maxOpen"`
	Loc             string `yaml:"loc" json:"loc"`
	MultiStatements bool   `yaml:"multiStatements" json:"multiStatements"`
	ParseTime       bool   `yaml:"parseTime" json:"parseTime"`
	ShowSql         bool   `yaml:"showSql" json:"showSql"`
}

func NewSimpleMysqlPool(setting *modules.SimpleSettingMysql) (*gorm.DB, error) {
	if setting == nil {
		return nil, fmt.Errorf("mysqlSetting is nil")
	}
	if setting.Username == "" {
		return nil, fmt.Errorf("lack of mysqlSetting.Username")
	}
	if setting.Password == "" {
		return nil, fmt.Errorf("lack of mysqlSetting.Password")
	}
	if setting.Host == "" {
		return nil, fmt.Errorf("lack of mysqlSetting.Host")
	}
	if setting.DBName == "" {
		return nil, fmt.Errorf("lack of mysqlSetting.DBName")
	}
	if setting.Charset == "" {
		return nil, fmt.Errorf("lack of mysqlSetting.Charset")
	}

	var buf bytes.Buffer
	buf.WriteString(setting.Username)
	buf.WriteString(":")
	buf.WriteString(setting.Password)
	buf.WriteString("@tcp(")
	buf.WriteString(setting.Host)
	buf.WriteString(")/")
	buf.WriteString(setting.DBName)
	buf.WriteString("?charset=")
	buf.WriteString(setting.Charset)
	buf.WriteString("&parseTime=" + strconv.FormatBool(setting.ParseTime))
	buf.WriteString("&multiStatements=" + strconv.FormatBool(setting.MultiStatements))

	if setting.Loc == "" {
		buf.WriteString("&loc=Local")
	} else {
		buf.WriteString("&loc=" + url.QueryEscape(setting.Loc))
	}

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

	return db, nil
}

/*
buf.WriteString(mysqlSetting.UserName)
	buf.WriteString(":")
	buf.WriteString(mysqlSetting.Password)
	buf.WriteString("@tcp(")
	buf.WriteString(mysqlSetting.Host)
	buf.WriteString(")/")
	buf.WriteString(mysqlSetting.DBName)
	buf.WriteString("?charset=")
	buf.WriteString(mysqlSetting.Charset)
	buf.WriteString("&parseTime=" + strconv.FormatBool(mysqlSetting.ParseTime))
	buf.WriteString("&multiStatements=" + strconv.FormatBool(mysqlSetting.MultiStatements))
	if mysqlSetting.ConnectionTimeout != "" {
		buf.WriteString(fmt.Sprintf("&timeout=%v", mysqlSetting.ConnectionTimeout))
	}
	if mysqlSetting.WriteTimeout != "" {
		buf.WriteString(fmt.Sprintf("&writeTimeout=%v", mysqlSetting.WriteTimeout))
	}
	if mysqlSetting.ReadTimeout != "" {
		buf.WriteString(fmt.Sprintf("&readTimeout=%v", mysqlSetting.ReadTimeout))
	}
	if mysqlSetting.Loc == "" {
		buf.WriteString("&loc=Local")
	} else {
		buf.WriteString("&loc=" + url.QueryEscape(mysqlSetting.Loc))
	}
*/
