package bridge

import (
	"context"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"keeper/app/pkg/serializer"
	"sync"
	"time"
)

type ConnPool struct {
	Conid  string
	Status map[string]string
}

var lock sync.RWMutex

type ServerConnections struct {
	Ctx        context.Context
	Opened     []*OpenedItem
	LastPinged map[string]UnixTime
}

type PingRequest struct {
	Connections []string
}

type OpenedItem struct {
	Conid        string
	Connection   map[string]interface{}
	Status       map[string]string
	Disconnected bool
}

func NewServerConnections() *ServerConnections {
	return &ServerConnections{
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

	newOpened := &OpenedItem{
		Conid:        conid,
		Connection:   connection,
		Status:       map[string]string{"name": "pending"},
		Disconnected: false,
	}

	sc.Opened = append(sc.Opened, newOpened)
	runtime.EventsEmit(Application.ctx, "server-status-changed")
	//socket.emitChanged(`server-status-changed`);
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
