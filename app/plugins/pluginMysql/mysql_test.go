package pluginMysql

import (
	"fmt"
	"testing"
)

func TestConnect(t *testing.T) {
	connect, err := Connect("127.0.0.1", "3306", "root", "123456", "mysql")
	if err != nil {
		return
	}
	query, err := connect.Query(`select CONSTRAINT_NAME as constraintName
    from information_schema.TABLE_CONSTRAINTS
    where CONSTRAINT_SCHEMA = 'yami_shops' and constraint_type = 'UNIQUE'`)
	var constraintName string
	if err != nil {
		return
	}

	for query.Next() {
		query.Scan(&constraintName)
		fmt.Printf("%s", constraintName)
	}
}
