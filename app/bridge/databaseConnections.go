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
	Opened             []map[string]interface{}
	Closed             map[string]interface{}
	DatabaseConnection *sideQuests.DatabaseConnectionHandlers
	Ch                 chan *modules.EchoMessage
}

func NewDatabaseConnections() *DatabaseConnections {
	ch := make(chan *modules.EchoMessage)
	return &DatabaseConnections{
		Ch:                 ch,
		DatabaseConnection: sideQuests.NewDatabaseConnectionHandlers(ch),
	}
}

func (dc *DatabaseConnections) Refresh(req *DatabaseKeepOpenRequest) *serializer.Response {
	if !req.KeepOpen {
		dc.close(req.Conid, req.Database, true)
	}
	dc.ensureOpened(req.Conid, req.Database)
	return serializer.SuccessData("", map[string]string{"status": "ok"})
}

func (dc *DatabaseConnections) SyncModel(req *DatabaseRequest) *serializer.Response {
	dc.ensureOpened(req.Conid, req.Database)

	return serializer.SuccessData("", map[string]string{"status": "ok"})
}

type databaseConnections struct {
	Conid    string `json:"conid"`
	Database string `json:"database"`
}

type DatabaseRequest struct {
	databaseConnections
}

type DatabaseKeepOpenRequest struct {
	databaseConnections
	KeepOpen bool `json:"keepOpen"`
}

func (dc *DatabaseConnections) handleStructure(conid, database string, structure interface{}) {
	logger.Infof("structure handleStructure %s", tools.ToJsonStr(structure))

	existing, ok := lo.Find[map[string]interface{}](dc.Opened, func(item map[string]interface{}) bool {
		return item["conid"] == conid && item["database"] == database
	})

	if existing == nil || !ok {
		return
	}

	existing["structure"] = structure

	runtime.EventsEmit(Application.ctx, "database-structure-changed", &databaseConnections{
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
		dc.DatabaseConnection.Ping()
	} else {
		existing = dc.ensureOpened(req.Conid, req.Database)
	}

	res := map[string]interface{}{"status": "ok"}
	if existing != nil {
		res["connectionStatus"] = existing["status"]
	} else {
		res["connectionStatus"] = nil
	}
	return serializer.SuccessData(serializer.SUCCESS, res)
}

//{"conid":"c31f5609-eb18-4d96-9fea-a29dc52f1c1d","database":"admin"}
func (dc *DatabaseConnections) ensureOpened(conid, database string) map[string]interface{} {
	existing, ok := lo.Find[map[string]interface{}](dc.Opened, func(item map[string]interface{}) bool {
		return item["conid"] == conid && item["database"] == database
	})

	if existing != nil && ok {
		logger.Infof("123456 :%s", tools.ToJsonStr(existing))
		return existing
	}

	connection := getCore(conid, false)
	lastClosed, ok := dc.Closed[fmt.Sprintf("%s/%s", conid, database)]

	newOpened := map[string]interface{}{
		"conid":         conid,
		"database":      database,
		"structure":     nil,
		"serverVersion": nil,
		"connection":    connection,
		"status":        &OpenedStatus{Name: "pending"},
	}

	if lastClosed == nil || !ok {
		newOpened["structure"] = schema.CreateEmptyStructure()
	} else {
		logger.Infof("newOpened.Opened : %s", tools.ToJsonStr(dc.Opened))
		logger.Infof("newOpened.Closed : %s", tools.ToJsonStr(dc.Closed))
	}

	dc.Opened = append(dc.Opened, newOpened)
	go dc.DatabaseConnection.Connect(
		map[string]interface{}{
			"conid":      conid,
			"database":   database,
			"connection": lo.Assign[string, interface{}](connection, map[string]interface{}{"database": database}),
		}, nil)
	go dc.listener(conid, database, dc.Ch)

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
	for {
		message, ok := <-chData
		if message != nil {
			switch message.MsgType {
			case "structure":
				dc.handleStructure(conid, database, message.Payload)
			case "structureTime":
				dc.handleStructureTime(conid, database, message.Payload.(code.UnixTime))
			case "exit":
				dc.close(conid, database, false)
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
	databaseConnections
	Kill bool `json:"kill"`
}

func (dc *DatabaseConnections) close(conid, database string, kill bool) {
	existing, ok := lo.Find[map[string]interface{}](dc.Opened, func(item map[string]interface{}) bool {
		if item[conidkey] != nil && item[conidkey].(string) == conid && item["database"] == database {
			return true
		} else {
			return false
		}
	})

	if existing != nil && ok {
		existing["disconnected"] = true
		dc.Opened = lo.Filter[map[string]interface{}](dc.Opened, func(obj map[string]interface{}, _ int) bool {
			return obj[conidkey].(string) != conid || obj["database"].(string) != database
		})
		dc.Closed[fmt.Sprintf("%s/%s", conid, database)] = map[string]interface{}{
			"status": existing["status"],
			"name":   "error",
		}

		runtime.EventsEmit(Application.ctx, "database-status-changed", &databaseConnections{
			Conid:    conid,
			Database: database,
		})
	}
}

func (dc *DatabaseConnections) closeAll(conid string, kill bool) {
	list := lo.Filter[map[string]interface{}](dc.Opened, func(obj map[string]interface{}, _ int) bool {
		return obj[conidkey].(string) != conid
	})

	for _, v := range list {
		dc.close(conid, v["database"].(string), kill)
	}
}

func (dc *DatabaseConnections) Disconnect(req *DatabaseRequest) *serializer.Response {
	dc.close(req.Conid, req.Database, true)
	return serializer.SuccessData("", &OpenedStatus{Name: "ok"})
}
