package bridge

import (
	"errors"
	"fmt"
	"github.com/samber/lo"
	"keeper/app/code"
	"keeper/app/modules"
	"keeper/app/pkg/logger"
	"keeper/app/pkg/serializer"
	"keeper/app/pkg/standard"
	"keeper/app/sideQuests"
	"keeper/app/tools"
	"keeper/app/utility"
	"sync"
)

var lock sync.RWMutex

const conidkey = "conid"

type ServerConnections struct {
	PoolMap                 map[string]standard.SqlStandard
	Closed                  map[string]interface{}
	Opened                  []map[string]interface{}
	LastPinged              map[string]code.UnixTime
	ServerConnectionChannel *sideQuests.ServerConnection
	ch                      chan *modules.EchoMessage
}

func NewServerConnections() *ServerConnections {
	ch := make(chan *modules.EchoMessage)
	return &ServerConnections{
		PoolMap:                 make(map[string]standard.SqlStandard),
		Closed:                  make(map[string]interface{}),
		LastPinged:              make(map[string]code.UnixTime),
		ServerConnectionChannel: sideQuests.NewServerConnection(ch),
		ch:                      ch,
	}
}

func (sc *ServerConnections) handleDatabases(conid string, databases interface{}) {
	existing, ok := lo.Find[map[string]interface{}](sc.Opened, func(item map[string]interface{}) bool {
		return item[conidkey] != nil && item[conidkey].(string) == conid
	})

	if existing == nil || !ok {
		return
	}

	existing["databases"] = databases
	utility.EmitChanged(Application.ctx, fmt.Sprintf("database-list-changed-%s", conid))
}

func (sc *ServerConnections) handleVersion(conid, version string) {

}

func (sc *ServerConnections) handleStatus(conid string, status map[string]string) {
	existing, ok := lo.Find[map[string]interface{}](sc.Opened, func(item map[string]interface{}) bool {
		return item[conidkey] != nil && item[conidkey].(string) == conid
	})

	if existing == nil || !ok {
		return
	}

	existing["status"] = &modules.OpenedStatus{Name: status["name"], Message: status["message"]}
	utility.EmitChanged(Application.ctx, "server-status-changed")
}

func (sc *ServerConnections) handlePing() {}

func (sc *ServerConnections) ensureOpened(conid string) map[string]interface{} {
	lock.Lock()
	defer lock.Unlock()

	existing, ok := lo.Find[map[string]interface{}](sc.Opened, func(x map[string]interface{}) bool {
		uuid, ok := x[conidkey].(string)
		return ok && uuid == conid
	})

	defer utility.EmitChanged(Application.ctx, "server-status-changed")
	if existing != nil && ok {
		if err := sc.checker(conid); err != nil {
			logger.Info("existing [%s]", tools.ToJsonStr(existing))
			sc.Close(conid, true)
		}
		return existing
	}

	connection := getCore(conid, false)
	newOpened := map[string]interface{}{
		conidkey:       conid,
		"status":       &modules.OpenedStatus{Name: "pending"},
		"databases":    []interface{}{},
		"connection":   connection,
		"disconnected": false,
	}

	sc.Opened = append(sc.Opened, newOpened)

	if sc.Closed != nil {
		delete(sc.Closed, conid)
	}

	go sc.ServerConnectionChannel.Connect(sc.ch, connection)
	go sc.pipeHandler(newOpened, sc.ch)

	return newOpened
}

func (sc *ServerConnections) checker(conid string) error {
	if sc.PoolMap[conid] != nil {
		return sc.PoolMap[conid].Ping()
	}

	return errors.New("invalid memory address or nil pointer dereference")
}

func (sc *ServerConnections) ListDatabases(request map[string]string) *serializer.Response {
	if request == nil || request[conidkey] == "" {
		return serializer.Fail(serializer.IdNotEmpty)
	}
	opened := sc.ensureOpened(request[conidkey])
	return serializer.SuccessData(serializer.SUCCESS, opened["databases"])
}

func (sc *ServerConnections) ServerStatus() interface{} {
	values := map[string]interface{}{}
	for key, val := range sc.Closed {
		values[key] = val
	}
	for _, driver := range sc.Opened {
		statusObj, ok := driver["status"].(*modules.OpenedStatus)
		if ok && statusObj != nil {
			values[driver[conidkey].(string)] = statusObj
		}
	}

	return serializer.SuccessData("", values)
}

func (sc *ServerConnections) Ping(connections []string) *serializer.Response {
	for _, conid := range lo.Uniq[string](connections) {
		last := sc.LastPinged[conid]
		if last > 0 && tools.NewUnixTime()-last < tools.GetUnixTime(30*1000) {
			//continue
		}

		sc.LastPinged[conid] = tools.NewUnixTime()
		sc.ensureOpened(conid)
		sc.ServerConnectionChannel.NewTime()
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
			sc.Opened = lo.Filter[map[string]interface{}](sc.Opened, func(x map[string]interface{}, _ int) bool {
				uuid, ok := x[conid].(string)
				return ok && uuid != conid
			})
			sc.Closed[conid] = map[string]interface{}{
				"name":   "error",
				"status": existing["status"].(*modules.OpenedStatus),
			}
			sc.LastPinged[conid] = 0
		}

		utility.EmitChanged(Application.ctx, "server-status-changed")
	}
}

type ServerRefreshRequest struct {
	Conid    string `json:"conid"`
	KeepOpen bool   `json:"keepOpen"`
}

func (sc *ServerConnections) Refresh(req *ServerRefreshRequest) *serializer.Response {
	if !req.KeepOpen {
		sc.Close(req.Conid, true)
	}
	sc.ensureOpened(req.Conid)
	return serializer.SuccessData("", map[string]string{
		"status": "ok",
	})
}

func (sc *ServerConnections) pipeHandler(newOpened map[string]interface{}, chData <-chan *modules.EchoMessage) {
	for {
		message, ok := <-chData
		logger.Infof("current: %s", message.MsgType)
		conid := newOpened[conidkey].(string)
		if message != nil {
			switch message.MsgType {
			case "status":
				sc.handleStatus(conid, message.Payload.(map[string]string))
			case "databases":
				sc.handleDatabases(conid, message.Payload)
			case "ping":
				sc.handlePing()
			case "exit":
				if !newOpened["disconnected"].(bool) {
					sc.Close(conid, true)
				}
			case "pool":
				sc.PoolMap[conid] = message.Payload.(standard.SqlStandard)
			}

			if !ok {
				if !newOpened["disconnected"].(bool) {
					sc.Close(conid, true)
				}
				break
			}
		}
	}
}
