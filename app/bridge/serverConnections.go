package bridge

import (
	"context"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"keeper/app/pkg/serializer"
	"keeper/app/tools"
	"sync"
	"time"
)

var lock sync.RWMutex

type ServerConnections struct {
	Ctx        context.Context
	Closed     map[string]string
	Opened     []map[string]interface{}
	LastPinged map[string]UnixTime
}

type PingRequest struct {
	Connections []string
}

type OpenedItem struct {
	Conid        string
	Connection   map[string]interface{}
	Disconnected bool
}

type status struct {
	name string
}

func NewServerConnections() *ServerConnections {
	return &ServerConnections{
		Closed:     make(map[string]string),
		LastPinged: make(map[string]UnixTime),
	}
}

func (sc *ServerConnections) ListDatabases(request map[string]string) interface{} {
	if request["conid"] == "" {
		return serializer.Fail(Application.ctx, "")
	}

	//https://esc.show/article/Golang-GUI-kai-fa-zhi-Webview

	return nil
}

func (sc *ServerConnections) getCore(conid string, mask bool) map[string]interface{} {
	if conid == "" {
		return nil
	}

	return JsonLinesDatabase.Get(conid)
}

func (sc *ServerConnections) ensureOpened(conid string) {
	lock.Lock()
	defer lock.Unlock()
	//var existing bool
	//for _, x := range sc.Opened {
	//	if x != nil && x[conid] != "" {
	//		existing = true
	//		break
	//	}
	//}
	//
	//if existing {
	//	return
	//}

	connection := sc.getCore(conid, false)
	newOpened := tools.MergeMaps(connection, map[string]interface{}{
		"conid":        conid,
		"status":       map[string]string{"name": "pending"},
		"disconnected": false,
	})

	sc.Opened = append(sc.Opened, newOpened)
	if sc.Closed != nil && sc.Closed[conid] != "" {
		delete(sc.Closed, conid)
	}

	runtime.EventsEmit(Application.ctx, "server-status-changed")
}

func (sc *ServerConnections) ServerStatus() interface{} {
	return serializer.SuccessData(Application.ctx, "", map[string]status{
		"efdc46d9-fed2-43d7-b506-53514b0a2559": {name: "ok"},
		"de5bb0d8-2a7c-4de6-92db-b60606a83c93": {name: "pending"},
	})
}

func (sc *ServerConnections) Ping(request *PingRequest) interface{} {
	//_.uniq(connections)

	if request == nil {
		return serializer.Fail(context.Background(), "")
	}

	for _, conid := range request.Connections {
		last := sc.LastPinged[conid]
		if last != 0 && UnixTime(time.Now().Unix())-last < UnixTime(30*1000) {
			//return Promise.resolve();
			return serializer.SuccessData(Application.ctx, "", map[string]string{"status": "ok"})
		}

		sc.LastPinged[conid] = UnixTime(time.Now().Unix())
		sc.ensureOpened(conid)
	}

	return serializer.SuccessData(Application.ctx, "", map[string]string{"status": "ok"})
}

func (sc *ServerConnections) Close(conid string, kill bool) {

}

func (sc *ServerConnections) Refresh(conid string) interface{} {
	sc.Close(conid, true)

	sc.ensureOpened(conid)

	return serializer.SuccessData(Application.ctx, "", map[string]string{
		"status": "ok",
	})
}

func (sc *ServerConnections) handleDatabases(conid, databases string) {
	var existing map[string]interface{}
	for _, x := range sc.Opened {
		if id, ok := x["conid"]; ok && id != nil && id.(string) == conid {
			existing = x
			break
		}
	}

	if existing == nil {
		return
	}
}

func (sc *ServerConnections) handleVersion(conid, version string) {

}

func (sc *ServerConnections) handleStatus(conid string, status map[string]string) {
	var existing map[string]interface{}
	for _, x := range sc.Opened {
		if id, ok := x["conid"]; ok && id != nil && id.(string) == conid {
			existing = x
			break
		}
	}

	if existing == nil {
		return
	}
	existing["status"] = status
	runtime.EventsEmit(Application.ctx, "server-status-changed")
}

func (sc *ServerConnections) handlePing() {

}
