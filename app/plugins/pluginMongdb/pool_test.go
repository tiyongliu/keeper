package pluginMongdb

import (
	"fmt"
	"keeper/app/pkg/logger"
	"keeper/app/plugins/modules"
	"keeper/app/utility"
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

	fmt.Println(utility.ToJsonStr(databases))
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
	tables, err := pool.(*MongoDBDrivers).Collections("local")
	if err == nil {
		logger.Infof("list %s", utility.ToJsonStr(tables))
	}
}
