package plugin_mondb

import (
	"fmt"
	"keeper/app/modules"
	"keeper/app/tools"
	"testing"
)

func TestListDatabases(t *testing.T) {
	pool, err := NewSimpleMongoDBPool(&modules.SimpleSettingMongoDB{
		Host:     "localhost",
		Port:     "27017",
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
