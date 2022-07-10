package plugin_mongdb

import (
	"fmt"
	"keeper/app/modules"
	"keeper/app/pkg/logger"
	"keeper/app/tools"
	"testing"
)

func TestListDatabases(t *testing.T) {
	pool, err := NewSimpleMongoDBPool(&modules.SimpleSettingMongoDB{
		Host: "localhost",
		Port: "27017",
	})

	if err != nil {
		fmt.Printf("err: %v \n", err)
		return
	}

	databases, err := pool.ListDatabases()
	if err != nil {
		fmt.Printf("err: %v \n", err)
		return
	}

	fmt.Println(tools.ToJsonStr(databases))
}

func TestTables(t *testing.T) {
	pool, err := NewSimpleMongoDBPool(&modules.SimpleSettingMongoDB{
		Host: "localhost",
		Port: "27017",
	})

	if err != nil {
		fmt.Printf("err: %v \n", err)
		return
	}
	tables, err := pool.Tables()
	if err == nil {
		logger.Infof("list %s", tools.ToJsonStr(tables))
	}
}
