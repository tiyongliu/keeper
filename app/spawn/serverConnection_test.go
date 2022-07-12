package spawn

import (
	"fmt"
	"testing"
)

func TestConnectUtility(t *testing.T) {
	connection := connectUtility(map[string]interface{}{
		"_id":      "65685eb0-4082-4363-903f-2e4742df3da8",
		"engine":   "mysql",
		"host":     "localhost",
		"password": "crypt:3ba94bc76d3cd137f665c5839c16bbb5e7d77316da557541a36dcf28e0263c80X7kLrbN8WCG22VUmpBqVBFNhbHRlZF9fwfTLx/3iV8ewqNtTOn/7wOAE4QGaGQ9s",
		"port":     "3306",
		"username": "root",
	})

	fmt.Println(connection)
}
