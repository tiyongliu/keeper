package plugin_mondb

import (
	"dbbox/app/src/modules"
	"fmt"
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

	version, err := pool.GetVersion()
	if err != nil {
		fmt.Printf("err: %v \n", err)
		return
	}

	fmt.Println(version)
}
