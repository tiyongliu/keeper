package sideQuests

import (
	"fmt"
	"github.com/samber/lo"
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

func TestContains(t *testing.T) {
	fmt.Println(lo.Contains[string]([]string{"0", "1", "2", "3", "4", "5"}, "1"))
}

var b []map[string]interface{}

func TestKeyBy(t *testing.T) {
	b = append(b, map[string]interface{}{
		"connid": "123",
	})

	b = append(b, map[string]interface{}{
		"connid": "456",
	})

	by := lo.KeyBy[string, map[string]interface{}](b, func(m map[string]interface{}) string {
		return "conid"
	})
	fmt.Println(by)
}

func TestMapValues(t *testing.T) {
	m1 := map[string]interface{}{
		"_id":    "75f6c2d7-65fd-4d8f-afa1-8cd615ee153b",
		"engine": "mongo",
		"host":   "localhost",
		"port":   "27017",
	}

	m2 := lo.MapValues[string, interface{}, string](m1, func(x interface{}, r string) string {
		fmt.Println(x, r)
		return "status"
	})

	fmt.Println(m2)
}
