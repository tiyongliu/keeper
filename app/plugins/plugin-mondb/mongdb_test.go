package plugin_mondb

import (
	"fmt"
	"keeper/app/internal/tools"
	"keeper/app/modules"
	"testing"
)

func TestConnect(t *testing.T) {
	Connect()
}

func Test_Pool(t *testing.T) {
	pool, err := NewSimpleMongoDBPool(&modules.SimpleSettingMongoDB{
		Host: "localhost",
		Port: "27017",
	})

	if err != nil {
		fmt.Printf("err: %v \n", err)
		return
	}

	defer pool.Close()

	version, err := pool.GetVersion()
	if err != nil {
		fmt.Printf("err: %v \n", err)
		return
	}

	fmt.Println(version)

	databases, err := pool.ListDatabases()
	if err != nil {
		fmt.Printf("err: %v \n", err)
		return
	}

	fmt.Printf("list: %s", tools.ToJsonStr(databases))

}
