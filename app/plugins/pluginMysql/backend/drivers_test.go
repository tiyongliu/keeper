package backend

import (
	"fmt"
	"keeper/app/db/adapter"
	"testing"
)

func TestNewAnalyser(t *testing.T) {
	driver, err := adapter.NewCompatDriver().Open(map[string]interface{}{
		"username": "root",
		"password": "123456",
		"port":     "3306",
		"database": "",
		"host":     "localhost",
	})
	if err != nil {
		fmt.Printf("err: %v \n", err)
		return
	}
	NewAnalyser(driver, "yami_shops").RunAnalysis()
}
