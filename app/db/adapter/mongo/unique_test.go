package mongo

import (
	"encoding/json"
	"fmt"
	"keeper/app/db"
	"keeper/app/db/standard/modules"
	"testing"
)

func TestCountDocuments(t *testing.T) {
	getDevice(func(session db.Session) {
		driver, ok := session.(*Source)
		if !ok && driver == nil {
			return
		}
		var condition map[string]interface{}
		sql := `{"$and":[{"$or":[{"pid":{"$regex":".*5284.*","$options":"i"}},{"pid":5284}]}]}`
		//	sql := `{"$and":[{"_id":{"$regex":".*SC-201907201119-1674316358023.*","$options":"i"}}]}`
		err := json.Unmarshal([]byte(sql), &condition)
		if err != nil {
			return
		}

		documents, err := driver.ReadCollection("local", &modules.CollectionDataOptions{
			PureName:       "startup_log",
			CountDocuments: true,
			Condition:      condition,
		})
		fmt.Println(err)
		fmt.Println(documents)
	})
}
func TestAggregate(t *testing.T) {
	getDevice(func(session db.Session) {
		driver, ok := session.(*Source)
		if !ok && driver == nil {
			return
		}
		documents, err := driver.ReadCollection("local", &modules.CollectionDataOptions{
			PureName:       "startup_log",
			CountDocuments: true,
			Limit:          50,
			Skip:           0,
			Aggregate:      nil,
		})
		fmt.Println(err)
		fmt.Println(documents)
	})
}
