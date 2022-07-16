package bridge

import (
	"fmt"
	"keeper/app/code"
	"keeper/app/modules"
	"keeper/app/pkg/logger"
	"keeper/app/pkg/serializer"
	"keeper/app/schema"
	"keeper/app/sideQuests"
	"keeper/app/tools"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/samber/lo"
)

const databaseKey = "database"

type DatabaseConnections struct {
	Opened []map[string]interface{}
	Closed map[string]interface{}
}

func NewDatabaseConnections() *DatabaseConnections {
	return &DatabaseConnections{}
}

func (dc *DatabaseConnections) Refresh(conid string) {

}

type databaseConnectionsBase struct {
	Conid    string `json:"conid"`
	Database string `json:"database"`
}

type DatabaseRequest struct {
	databaseConnectionsBase
}

func (dc *DatabaseConnections) handleStructure(conid, database string, structure interface{}) {
	logger.Infof("structure handleStructure %s", tools.ToJsonStr(structure))

	existing, ok := lo.Find[map[string]interface{}](dc.Opened, func(item map[string]interface{}) bool {
		if item[conidkey] != nil && item[conidkey].(string) == conid && item["database"] == database {
			return true
		} else {
			return false
		}
	})

	if existing == nil || !ok {
		return
	}

	logger.Infof("structure EventsEmit %s", tools.ToJsonStr(structure))
	existing["structure"] = structure

	//runtime.EventsEmit(Application.ctx, fmt.Sprintf("database-structure-changed-%s-%s", conid, database))
	runtime.EventsEmit(Application.ctx, "database-structure-changed", &databaseConnectionsBase{
		Conid:    conid,
		Database: database,
	})
}

func (dc *DatabaseConnections) handleStructureTime(conid, database string, analysedTime code.UnixTime) {
	existing, ok := lo.Find[map[string]interface{}](dc.Opened, func(item map[string]interface{}) bool {
		if item[conidkey] != nil && item[conidkey].(string) == conid && item["database"] == database {
			return true
		} else {
			return false
		}
	})

	if existing == nil || !ok {
		return
	}

	existing["analysedTime"] = analysedTime

	runtime.EventsEmit(Application.ctx, fmt.Sprintf("database-status-changed-%s-%s", conid, database))
}

func (dc *DatabaseConnections) Ping(req *DatabaseRequest) *serializer.Response {
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
		if item[conidkey] != nil && item[conidkey].(string) == conid && item["database"] == database {
			return true
		} else {
			return false
		}
	})

	if existing != nil && ok {
		logger.Infof("123456 :%s", tools.ToJsonStr(existing))
		return existing
	}

	connection := getCore(conid, false)
	lastClosed, ok := dc.Closed[fmt.Sprintf("%s/%s", conid, database)]

	newOpened := map[string]interface{}{
		"conid":         conid,
		"serverVersion": nil,
		"status":        &OpenedStatus{Name: "pending"},
		"structure":     nil,
	}

	if lastClosed == nil || !ok {
		newOpened["structure"] = schema.CreateEmptyStructure()
	} else {
		logger.Infof("newOpened.Opened : %s", tools.ToJsonStr(dc.Opened))
		logger.Infof("newOpened.Closed : %s", tools.ToJsonStr(dc.Closed))
	}

	dc.Opened = append(dc.Opened, newOpened)

	ch := make(chan *modules.EchoMessage)

	go sideQuests.NewDatabaseConnectionHandlers(ch).Connect(
		map[string]interface{}{
			"conid":      conid,
			"database":   database,
			"connection": lo.Assign[string, interface{}](connection, map[string]interface{}{"database": database}),
		}, nil)
	go dc.listener(conid, database, ch)

	return newOpened
}

//{"conid":"d0fc6ec0-fae2-11ec-ad02-b72b9a6655f8","database":"erd"}
func (dc *DatabaseConnections) Structure(req *DatabaseRequest) interface{} {
	if req.Conid == "__model" {
		//todo  const model = await importDbModel(database);
	}

	opened := dc.ensureOpened(req.Conid, req.Database)
	return opened["structure"]
}

func (dc *DatabaseConnections) listener(conid, database string, chData <-chan *modules.EchoMessage) {
	//structure

	for {
		message, ok := <-chData
		if message != nil {
			switch message.MsgType {
			case "structure":
				dc.handleStructure(conid, database, message.Payload)
			case "structureTime":
				dc.handleStructureTime(conid, database, message.Payload.(code.UnixTime))
			}
		}
		if !ok {
			break
		}
	}
}

func (dc *DatabaseConnections) ServerVersion(req *DatabaseRequest) *serializer.Response {
	if req == nil || req.Conid == "" {
		return serializer.Fail(serializer.ParamsErr)
	}

	opened := dc.ensureOpened(req.Conid, req.Database)
	return serializer.SuccessData("", opened["serverVersion"])
}

func (dc *DatabaseConnections) Status(req *DatabaseRequest) *serializer.Response {
	existing, ok := lo.Find[map[string]interface{}](dc.Opened, func(item map[string]interface{}) bool {
		if item[conidkey] != nil && item[conidkey].(string) == req.Conid && item["database"] == req.Database {
			return true
		} else {
			return false
		}
	})

	if existing != nil && ok {
		serializer.SuccessData("", map[string]interface{}{

			"analysedTime": existing["analysedTime"],
		})
	}
	lastClosed, ok := dc.Closed[fmt.Sprintf("%s/%s", req.Conid, req.Database)]
	if lastClosed != nil && ok {
		return serializer.SuccessData("", map[string]interface{}{
			"analysedTime": lastClosed.(map[string]interface{})["analysedTime"],
		})
	}

	return serializer.Fail("Not connected")
}

func (dc *DatabaseConnections) sendRequest() {

}

type DatabaseKillRequest struct {
	databaseConnectionsBase
	Kill bool `json:"kill"`
}

func (dc *DatabaseConnections) Close(req *DatabaseKillRequest) {
	existing, ok := lo.Find[map[string]interface{}](dc.Opened, func(item map[string]interface{}) bool {
		if item[conidkey] != nil && item[conidkey].(string) == req.Conid && item["database"] == req.Database {
			return true
		} else {
			return false
		}
	})

	if existing != nil && ok {
		existing["disconnected"] = true
		dc.Opened = lo.Filter[map[string]interface{}](dc.Opened, func(obj map[string]interface{}, _ int) bool {
			return obj[conidkey].(string) != req.Conid || obj["database"].(string) != req.Database
		})
		dc.Closed[fmt.Sprintf("%s/%s", req.Conid, req.Database)] = map[string]interface{}{
			"status": existing["status"],
			"name":   "error",
		}

		runtime.EventsEmit(Application.ctx, "database-status-changed", &databaseConnectionsBase{
			Conid:    req.Conid,
			Database: req.Database,
		})
	}
}

type CloseAllRequest struct {
	Conid string `json:"conid"`
	Kill  bool   `json:"kill"`
}

func (dc *DatabaseConnections) CloseAll(req *CloseAllRequest) {

}
