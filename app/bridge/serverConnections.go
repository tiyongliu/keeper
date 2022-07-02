package bridge

import (
	"github.com/samber/lo"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"keeper/app/code"
	"keeper/app/modules"
	"keeper/app/pkg/logger"
	"keeper/app/pkg/serializer"
	"keeper/app/sideQuests"
	"keeper/app/tools"
	"sync"
	"time"
)

var lock sync.RWMutex

const conidkey = "conid"

type ServerConnections struct {
	Closed     map[string]interface{}
	Opened     []map[string]interface{}
	LastPinged map[string]code.UnixTime
}

type OpenedStatus struct {
	Name string `json:"name"`
}

func NewServerConnections() *ServerConnections {
	return &ServerConnections{
		Closed:     make(map[string]interface{}),
		LastPinged: make(map[string]code.UnixTime),
	}
}

//https://esc.show/article/Golang-GUI-kai-fa-zhi-Webview
func (sc *ServerConnections) ListDatabases(request string) interface{} {
	if request == "" {
		return serializer.Fail(Application.ctx, "")
	}

	opened := sc.ensureOpened(request)

	logger.Infof("opened: %s", tools.ToJsonStr(opened))
	return nil
}

func (sc *ServerConnections) getCore(conid string, mask bool) map[string]interface{} {
	if conid == "" {
		return nil
	}

	return JsonLinesDatabase.Get(conid)
}

func (sc *ServerConnections) ensureOpened(conid string) map[string]interface{} {
	lock.Lock()
	defer lock.Unlock()

	existing, ok := lo.Find[map[string]interface{}](sc.Opened, func(item map[string]interface{}) bool {
		if item[conidkey].(string) == conid {
			return true
		} else {
			return false
		}
	})

	if existing != nil && ok {
		return existing
	}

	connection := sc.getCore(conid, false)

	newOpened := map[string]interface{}{
		conidkey:       conid,
		"status":       &OpenedStatus{Name: "pending"},
		"databases":    []interface{}{},
		"connection":   connection,
		"disconnected": false,
	}

	sc.Opened = append(sc.Opened, newOpened)

	if sc.Closed != nil && sc.Closed[conid] != "" {
		delete(sc.Closed, conid)
	}

	ch := make(chan *modules.EchoMessage)
	go sideQuests.NewMessageDriverHandlers(ch).Connect(connection)
	go sc.Listener(conid, ch)

	runtime.EventsEmit(Application.ctx, "server-status-changed")

	return newOpened
}

func (sc *ServerConnections) ServerStatus() interface{} {
	values := map[string]interface{}{}
	for _, driver := range sc.Opened {
		statusObj, ok := driver["status"].(*OpenedStatus)
		if ok {
			values[driver[conidkey].(string)] = statusObj
		}
	}

	return serializer.SuccessData(Application.ctx, "", values)
}

func (sc *ServerConnections) Ping(connections []string) interface{} {
	for _, conid := range lo.Uniq[string](connections) {
		last := sc.LastPinged[conid]
		if last > 0 && code.UnixTime(time.Now().Unix())-last < code.UnixTime(30*1000) {
			continue
		}

		sc.LastPinged[conid] = code.UnixTime(time.Now().Unix())
		sc.ensureOpened(conid)
	}

	return serializer.SuccessData(Application.ctx, "", map[string]string{"status": "ok"})
}

func (sc *ServerConnections) Close(conid string, kill bool) {
	existing, ok := lo.Find[map[string]interface{}](sc.Opened, func(item map[string]interface{}) bool {
		if item[conidkey].(string) != "" && item[conidkey].(string) == conid {
			return true
		} else {
			return false
		}
	})

	if existing != nil && ok {
		existing["disconnected"] = true
		if kill {
			//existing.subprocess.kill()
			sc.Opened = lo.Filter[map[string]interface{}](sc.Opened, func(obj map[string]interface{}, _ int) bool {
				uuid, ok := obj[conid].(string)
				return ok && uuid != conid
			})

			//{"_id":"75f6c2d7-65fd-4d8f-afa1-8cd615ee153b","conid":"75f6c2d7-65fd-4d8f-afa1-8cd615ee153b","disconnected":false,"engine":"mongo","host":"localhost","port":"27017","status":{"name":"pending"}}
			sc.Closed = map[string]interface{}{
				"name":   "error",
				"status": existing["status"].(*OpenedStatus),
			}
		}
		runtime.EventsEmit(Application.ctx, "server-status-changed")
	}
}

type RefreshRequest struct {
	Conid    string `json:"conid"`
	KeepOpen bool   `json:"keepOpen"`
}

func (sc *ServerConnections) Refresh(req *RefreshRequest) interface{} {
	if !req.KeepOpen {
		sc.Close(req.Conid, true)
	}
	sc.ensureOpened(req.Conid)

	return serializer.SuccessData(Application.ctx, "", map[string]string{
		"status": "ok",
	})
}

func (sc *ServerConnections) handleDatabases(conid string, databases interface{}) {
	existing, ok := lo.Find[map[string]interface{}](sc.Opened, func(item map[string]interface{}) bool {
		if item[conidkey] != nil && item[conidkey].(string) == conid {
			return true
		} else {
			return false
		}
	})

	if existing == nil || !ok {
		return
	}

	existing["databases"] = databases

	runtime.EventsEmit(Application.ctx, "database-list-changed-"+conid)
}

func (sc *ServerConnections) handleVersion(conid, version string) {

}

func (sc *ServerConnections) handleStatus(conid string, status *sideQuests.StatusMessage) {
	existing, ok := lo.Find[map[string]interface{}](sc.Opened, func(item map[string]interface{}) bool {
		if item[conidkey] != nil && item[conidkey].(string) == conid {
			return true
		} else {
			return false
		}
	})

	if existing == nil || !ok {
		return
	}

	existing["status"] = &OpenedStatus{Name: status.Name}

	runtime.EventsEmit(Application.ctx, "server-status-changed")
}

func (sc *ServerConnections) handlePing() {}

func (sc *ServerConnections) Listener(conid string, chData <-chan *modules.EchoMessage) {
	for {
		message, ok := <-chData
		if message != nil {
			switch message.MsgType {
			case "status":
				sc.handleStatus(conid, message.Payload.(*sideQuests.StatusMessage))
			case "databases":
				sc.handleDatabases(conid, message.Payload)
			}
		}

		if !ok {
			break
		}
	}
}
