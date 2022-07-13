package bridge

import (
	"fmt"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"keeper/app/pkg/serializer"
	"keeper/app/spawn"

	"github.com/samber/lo"
)

const databaseKey = "database"

type DatabaseConnections struct {
	Opened []map[string]interface{}
}

func NewDatabaseConnections() *DatabaseConnections {
	return &DatabaseConnections{}
}

func (dc *DatabaseConnections) Refresh(conid string) {

}

type DatabasePingRequest struct {
	Conid    string `json:"conid"`
	Database string `json:"database"`
}

func (dc *DatabaseConnections) handleStructure(conid, database string) {
	existing, ok := lo.Find[map[string]interface{}](dc.Opened, func(item map[string]interface{}) bool {
		if item[conidkey] != nil && item[conidkey].(string) == conid {
			return true
		} else {
			return false
		}
	})

	if existing != nil && ok {

	}

	runtime.EventsEmit(Application.ctx, fmt.Sprintf("database-structure-changed-%s-%s", conid, database))
}

func (dc *DatabaseConnections) Ping(req *DatabasePingRequest) interface{} {
	if req == nil || req.Conid == "" {
		return serializer.Fail(serializer.IdNotEmpty)
	}

	existing, ok := lo.Find[map[string]interface{}](dc.Opened, func(item map[string]interface{}) bool {
		if item[conidkey] != nil && item[conidkey].(string) == req.Conid {
			return true
		} else {
			return false
		}
	})

	if existing != nil && ok {

	} else {

	}

	return serializer.SuccessData(serializer.SUCCESS, map[string]interface{}{
		"status": "ok",
	})
}

//{"conid":"11485e70-e41e-11ec-aad8-95f9fdd48a30","database":"admin"}
func (dc *DatabaseConnections) ensureOpened(conid, database string) map[string]interface{} {
	existing, ok := lo.Find[map[string]interface{}](dc.Opened, func(item map[string]interface{}) bool {
		if item[conidkey] != nil && item[conidkey].(string) == conid {
			return true
		} else {
			return false
		}
	})

	if existing != nil && ok {
		return existing
	}

	connection := getCore(conid, false)

	newOpened := map[string]interface{}{
		"conid":     conid,
		"status":    &OpenedStatus{Name: "pending"},
		"structure": nil,
	}

	dc.Opened = append(dc.Opened, newOpened)

	go spawn.NewDatabaseConnectionHandlers().Connect(map[string]interface{}{
		"connection": lo.Assign[string, interface{}](connection, map[string]interface{}{"database": database}),
	}, nil)

	return newOpened
}

//{"conid":"d0fc6ec0-fae2-11ec-ad02-b72b9a6655f8","database":"erd"}
func (dc *DatabaseConnections) Structure(conid, database string) interface{} {
	if conid == "__model" {
		//todo  const model = await importDbModel(database);
	}

	opened := dc.ensureOpened(conid, database)
	return opened["structure"]
}
