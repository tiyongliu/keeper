package plugin_mysql

import (
	"fmt"
	"keeper/app/modules"
	"keeper/app/tools"
	"testing"
)

func TestConnect(t *testing.T) {
	connect, err := Connect("127.0.0.1", "3306", "root", "123456", "mysql")
	fmt.Println(err)
	query, err := connect.Query("select version() FROM DUAL;")
	if err != nil {

	}
	fmt.Println(query)
}

func TestQueryAll(t *testing.T) {
	connect, err := Connect("127.0.0.1", "3306", "root", "123456", "mysql")
	fmt.Println(err)

	all, err := QueryAll(connect, "select version() as version FROM DUAL;")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(all)
}

func TestGetVersion(t *testing.T) {
	result, err := NewMysql().GetVersion()
	fmt.Println(err)
	t.Logf("%v", result)
}

func TestPool(t *testing.T) {
	pool, err := NewSimpleMysqlPool(&modules.SimpleSettingMysql{
		Host:     "localhost",
		Username: "root",
		Password: "123456",
		Port:     "3306",
	})

	if err != nil {
		fmt.Printf("err: %v \n", err)
	}

	defer pool.Close()

	version, err := pool.GetVersion()
	if err != nil {
		fmt.Printf("err: %v \n", err)
	}

	fmt.Println(version)
}

func TestListDatabases(t *testing.T) {
	pool, err := NewSimpleMysqlPool(&modules.SimpleSettingMysql{
		Host:     "localhost",
		Username: "root",
		Password: "123456",
		Port:     "3306",
	})

	if err != nil {
		fmt.Printf("err: %v \n", err)
	}

	defer pool.Close()

	lastDatabases, err := pool.ListDatabases()
	if err != nil {
		fmt.Printf("err: %v \n", err)
	}

	TransformListDatabases(lastDatabases.([]string))

	fmt.Println(tools.ToJsonStr(TransformListDatabases(lastDatabases.([]string))))
}

func TestTable(t *testing.T) {
	pool, err := NewSimpleMysqlPool(&modules.SimpleSettingMysql{
		Host:     "localhost",
		Username: "root",
		Password: "123456",
		Port:     "3306",
	})

	if err != nil {
		fmt.Printf("err: %v \n", err)
	}

	defer pool.Close()
	pool.Tables()
}
