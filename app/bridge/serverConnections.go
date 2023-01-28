package bridge

import (
	"fmt"
	"github.com/samber/lo"
	"keeper/app/db/standard/modules"
	"keeper/app/db/stash"
	"keeper/app/internal/schema"
	"keeper/app/pkg/serializer"
	"keeper/app/sideQuests"
	"keeper/app/utility"
	"sync"
)

var lock sync.RWMutex

const conidkey = "conid"

type ServerConnections struct {
	Closed                  map[string]interface{}
	Opened                  []*schema.OpenedServerConnection
	LastPinged              map[string]utility.UnixTime
	ServerConnectionChannel *sideQuests.ServerConnection
}

func NewServerConnections() *ServerConnections {
	return &ServerConnections{
		Closed:                  make(map[string]interface{}),
		LastPinged:              make(map[string]utility.UnixTime),
		ServerConnectionChannel: sideQuests.NewServerConnection(),
	}
}

func (sc *ServerConnections) handleDatabases(conid string, databases interface{}) {
	existing := findByServerConnection(sc.Opened, conid)

	if existing == nil {
		return
	}

	existing.Databases = databases

	utility.EmitChanged(Application.ctx, fmt.Sprintf("database-list-changed-%s", conid))
}

func (sc *ServerConnections) handleVersion(conid string, version *modules.Version) {
	existing := findByServerConnection(sc.Opened, conid)

	if existing == nil {
		return
	}

	existing.Version = version
	utility.EmitChanged(Application.ctx, fmt.Sprintf("server-version-changed-%s", conid))
}

func (sc *ServerConnections) handleStatus(conid string, status *schema.OpenedStatus) {
	existing := findByServerConnection(sc.Opened, conid)

	if existing == nil {
		return
	}

	existing.Status = status

	utility.EmitChanged(Application.ctx, "server-status-changed")
}

func (sc *ServerConnections) handlePing() {}

func (sc *ServerConnections) ensureOpened(conid string) *schema.OpenedServerConnection {
	lock.Lock()
	defer lock.Unlock()
	existing := findByServerConnection(sc.Opened, conid)

	if existing != nil {
		return existing
	}

	connection := getCore(conid, false)
	if connection == nil {
		return nil
	}

	newOpened := &schema.OpenedServerConnection{
		Conid:        conid,
		Status:       &schema.OpenedStatus{Name: "pending"},
		Databases:    nil,
		Connection:   connection,
		Disconnected: false,
		Version:      nil,
	}

	sc.Opened = append(sc.Opened, newOpened)

	if sc.Closed != nil {
		delete(sc.Closed, conid)
	}
	utility.EmitChanged(Application.ctx, "server-status-changed")

	ch := make(chan *schema.EchoMessage)
	defer func() {
		sc.ServerConnectionChannel.ResetVars()
		go sc.ServerConnectionChannel.Connect(ch, connection)
		go sc.receiver(ch, conid)
	}()

	return newOpened
}

func (sc *ServerConnections) ListDatabases(request map[string]string) *serializer.Response {
	if request == nil || request[conidkey] == "" {
		return serializer.Fail(serializer.IdNotEmpty)
	}
	opened := sc.ensureOpened(request[conidkey])
	if opened == nil {
		return serializer.SuccessData(serializer.SUCCESS, nil)
	}
	return serializer.SuccessData(serializer.SUCCESS, opened.Databases)
}

func (sc *ServerConnections) ServerStatus() interface{} {
	values := map[string]interface{}{}
	for _, driver := range sc.Opened {
		values[driver.Conid] = driver.Status
	}
	for key, val := range sc.Closed {
		values[key] = val
	}
	return serializer.SuccessData(serializer.SUCCESS, values)
}

type ServerPingRequest struct {
	Connections []string `json:"connections"`
}

func (sc *ServerConnections) Ping(req *ServerPingRequest) *serializer.Response {
	for _, conid := range lo.Uniq[string](req.Connections) {
		last := sc.LastPinged[conid]
		//if driver, err := stash.GetStorageSession().GetItem(conid); err == nil {
		//	if err = driver.Ping(); err != nil {
		//		sc.Close(conid, true)
		//		continue
		//	}
		//}

		if last > 0 && utility.NewUnixTime()-last < utility.GetUnixTime(30*1000) {
			continue
		}

		sc.LastPinged[conid] = utility.NewUnixTime()
		sc.ensureOpened(conid)
		sc.ServerConnectionChannel.Ping()
	}

	return serializer.SuccessData(serializer.SUCCESS, map[string]string{"status": "ok"})
}

func (sc *ServerConnections) Close(conid string, kill bool) {
	existing := findByServerConnection(sc.Opened, conid)
	if existing != nil {
		existing.Disconnected = true
		if kill {
		}
		sc.Opened = lo.Filter[*schema.OpenedServerConnection](sc.Opened, func(x *schema.OpenedServerConnection, _ int) bool {
			return x.Conid != conid
		})
		sc.Closed[conid] = map[string]interface{}{
			"name":   "error",
			"status": existing.Status,
		}
		sc.LastPinged[conid] = 0
		_ = stash.GetStorageSession().RemoveItem(conid)
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
	return serializer.SuccessData(serializer.SUCCESS, map[string]string{
		"status": "ok",
	})
}

func (sc *ServerConnections) receiver(chData <-chan *schema.EchoMessage, conid string) {
	for {
		message, ok := <-chData
		if message != nil {
			if message.Err != nil {
				if existing := findByServerConnection(sc.Opened, conid); existing != nil && !existing.Disconnected {
					sc.Close(conid, false)
				}
			}
			switch message.MsgType {
			case "status":
				sc.handleStatus(conid, message.Payload.(*schema.OpenedStatus))
			case "version":
				sc.handleVersion(conid, message.Payload.(*modules.Version))
			case "databases":
				sc.handleDatabases(conid, message.Payload)
			case "ping":
				sc.handlePing()
			}
		}
		if !ok {
			break
		}
	}
}

func findByServerConnection(s []*schema.OpenedServerConnection, conid string) *schema.OpenedServerConnection {
	existing, ok := lo.Find[*schema.OpenedServerConnection](s, func(x *schema.OpenedServerConnection) bool {
		return x.Conid == conid
	})

	if existing != nil && ok {
		return existing
	}
	return nil
}
