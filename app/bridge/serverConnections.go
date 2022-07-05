package bridge

import (
	"keeper/app/code"
	"keeper/app/modules"
	"keeper/app/pkg/logger"
	"keeper/app/pkg/serializer"
	"keeper/app/sideQuests"
	"keeper/app/tools"
	"sync"
	"time"

	"github.com/samber/lo"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

var lock sync.RWMutex

const conidkey = "conid"

type ServerConnections struct {
	Closed         map[string]interface{}
	Opened         []map[string]interface{}
	LastPinged     map[string]code.UnixTime
	ch             chan *modules.EchoMessage
	DriverHandlers *sideQuests.MessageDriverHandlers
}

type OpenedStatus struct {
	Name string `json:"name"`
}

func NewServerConnections() *ServerConnections {
	ch := make(chan *modules.EchoMessage)
	return &ServerConnections{
		Closed:         make(map[string]interface{}),
		LastPinged:     make(map[string]code.UnixTime),
		ch:             ch,
		DriverHandlers: sideQuests.NewMessageDriverHandlers(ch),
	}
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

	runtime.EventsEmit(Application.ctx, "database-list-changed", conid)
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

	runtime.EventsEmit(Application.ctx, "server-status-changed")

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

	if sc.Closed != nil {
		delete(sc.Closed, conid)
	}

	go sc.DriverHandlers.Connect(connection)
	go sc.consumer(newOpened, sc.ch)

	return newOpened
}

//https://esc.show/article/Golang-GUI-kai-fa-zhi-Webview
func (sc *ServerConnections) ListDatabases(request string) interface{} {
	if request == "" {
		return serializer.Fail(serializer.IdNotEmpty)
	}

	opened := sc.ensureOpened(request)

	logger.Infof("opened: %s", tools.ToJsonStr(opened))

	return serializer.SuccessData(serializer.SUCCESS, opened["databases"])
}

func (sc *ServerConnections) getCore(conid string, mask bool) map[string]interface{} {
	if conid == "" {
		return nil
	}

	return JsonLinesDatabase.Get(conid)
}

func (sc *ServerConnections) ServerStatus() interface{} {
	values := map[string]interface{}{}
	logger.Infof("data1111111111111111111111: %s", tools.ToJsonStr(sc.Opened))
	logger.Infof("Closed2222222: %s", tools.ToJsonStr(sc.Closed))

	for _, driver := range sc.Opened {
		statusObj, ok := driver["status"].(*OpenedStatus)
		if ok {
			values[driver[conidkey].(string)] = statusObj
		}

	}

	return serializer.SuccessData("", values)
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

	return serializer.SuccessData("", map[string]string{"status": "ok"})
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

	return serializer.SuccessData("", map[string]string{
		"status": "ok",
	})
}

//func (sc *ServerConnections) consumer(conid string, chData <-chan *modules.EchoMessage) {
func (sc *ServerConnections) consumer(newOpened map[string]interface{}, chData <-chan *modules.EchoMessage) {
	for {
		message, ok := <-chData

		logger.Infof("chan message -<: %s", tools.ToJsonStr(message))
		conid := newOpened[conidkey].(string)
		if message != nil {
			switch message.MsgType {
			case "status":
				sc.handleStatus(conid, message.Payload.(*sideQuests.StatusMessage))
			case "databases":
				sc.handleDatabases(conid, message.Payload)
			case "ping":
				sc.handlePing()
			case "exit":
				if !newOpened["disconnected"].(bool) {
					sc.Close(conid, true)
				}
			}
		}

		if !ok {
			break
		}
	}
}
