package plugin_mysql

import (
	"fmt"
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
