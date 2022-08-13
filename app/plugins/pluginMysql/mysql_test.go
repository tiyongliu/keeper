package pluginMysql

import (
	"fmt"
	"keeper/app/pkg/logger"
	"keeper/app/plugins/modules"
	"keeper/app/utility"
	"testing"
)

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

	fmt.Println(utility.ToJsonStr(TransformListDatabases(lastDatabases.([]string))))
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

	tables, err := pool.(*MysqlDrivers).Tables("yami_shops", "tz_user")
	if err == nil {
		logger.Infof("list %s", utility.ToJsonStr(tables))
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
		logger.Infof("list %s", utility.ToJsonStr(columns))
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
		logger.Infof("list %s", utility.ToJsonStr(keys))
	}
}

func TestForeignKeys(t *testing.T) {
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

	keys, err := driver.ForeignKeys("yami_shops", "tz_user")
	if err == nil {
		logger.Infof("list %s", utility.ToJsonStr(keys))
	}
}
