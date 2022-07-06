package bridge

import (
	"keeper/app/pkg/serializer"

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

	newOpened := map[string]interface{}{
		"conid":  conid,
		"status": &OpenedStatus{Name: "pending"},
	}

	dc.Opened = append(dc.Opened, newOpened)

	return nil
}
