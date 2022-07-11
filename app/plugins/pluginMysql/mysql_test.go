package pluginMysql

import (
	"fmt"
	"keeper/app/modules"
	"keeper/app/pkg/logger"
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

/*
[2022-07-10 20:41:55.197][localhost_3306][000014][MYSQL]
SELECT * FROM `kb-dms`.`user_info` LIMIT 0, 1000
Time: 0.017s

[2022-07-10 20:41:55.215][localhost_3306][000015][MYSQL]
SHOW COLUMNS FROM `kb-dms`.`user_info`
Time: 0.001s

[2022-07-10 20:41:55.217][localhost_3306][000014][MYSQL]
SHOW TABLE STATUS LIKE 'user_info'
Time: 0.001s

[2022-07-10 20:41:55.218][localhost_3306][000014][MYSQL]
SHOW CREATE TABLE `kb-dms`.`user_info`
Time: 0.000s


*/
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

func TestTables(t *testing.T) {
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
	tables, err := pool.Tables()
	if err == nil {
		logger.Infof("list %s", tools.ToJsonStr(tables))
	}
}

func TestColumns(t *testing.T) {
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
	columns, err := pool.Columns("shop_go", "tz_user")
	if err == nil {
		logger.Infof("list %s", tools.ToJsonStr(columns))
		for _, item := range columns.([]*modules.MysqlTableColumn) {
			fmt.Println(string(item.ColumnComment.([]byte)))
		}
	}
}

func TestPrimaryKeys(t *testing.T) {
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
	driver, ok := pool.(*MysqlDrivers)
	if !ok && driver == nil {
		return
	}

	keys, err := driver.PrimaryKeys("shop_go", "tz_user")
	if err == nil {
		logger.Infof("list %s", tools.ToJsonStr(keys))
	}
}
