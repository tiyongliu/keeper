package bridge

import (
	"fmt"
	"github.com/samber/lo"
	"keeper/app/internal"
	"keeper/app/pkg/containers"
	"keeper/app/pkg/serializer"
	"keeper/app/pkg/standard"
	"keeper/app/sideQuests"
	"keeper/app/utility"
	"sync"
)

var lock sync.RWMutex

const conidkey = "conid"

type ServerConnections struct {
	Closed                  map[string]interface{}
	Opened                  []*containers.OpenedServerConnection
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

func (sc *ServerConnections) handleVersion(conid string, version *standard.VersionMsg) {
	existing := findByServerConnection(sc.Opened, conid)

	if existing == nil {
		return
	}

	existing.Version = version
	utility.EmitChanged(Application.ctx, fmt.Sprintf("server-version-changed-%s", conid))
}

func (sc *ServerConnections) handleStatus(conid string, status *containers.OpenedStatus) {
	existing := findByServerConnection(sc.Opened, conid)

	if existing == nil {
		return
	}

	existing.Status = status
	utility.EmitChanged(Application.ctx, "server-status-changed")
}

func (sc *ServerConnections) handlePing() {}

func (sc *ServerConnections) ensureOpened(conid string) *containers.OpenedServerConnection {
	lock.Lock()
	defer lock.Unlock()
	existing := findByServerConnection(sc.Opened, conid)

	if existing != nil {
		utility.EmitChanged(Application.ctx, "server-status-changed")
		return existing
	}

	connection := getCore(conid, false)
	newOpened := &containers.OpenedServerConnection{
		Conid:        conid,
		Status:       &containers.OpenedStatus{Name: "pending"},
		Databases:    nil,
		Connection:   connection,
		Disconnected: false,
		Version:      nil,
	}

	sc.Opened = append(sc.Opened, newOpened)

	if sc.Closed != nil {
		delete(sc.Closed, conid)
	}

	ch := make(chan *containers.EchoMessage)
	utility.EmitChanged(Application.ctx, "server-status-changed")

	defer func() {
		go sc.ServerConnectionChannel.Connect(ch, conid, connection)
		go sc.pipeHandler(ch, conid)
	}()

	return newOpened
}

func (sc *ServerConnections) checker(conid string) error {
	pool, err := internal.GetStoragePool(conid)
	if err != nil {
		return err
	}
	return pool.Ping()
}

func (sc *ServerConnections) ListDatabases(request map[string]string) *serializer.Response {
	if request == nil || request[conidkey] == "" {
		return serializer.Fail(serializer.IdNotEmpty)
	}
	opened := sc.ensureOpened(request[conidkey])
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

func (sc *ServerConnections) Ping(connections []string) *serializer.Response {
	for _, conid := range lo.Uniq[string](connections) {
		last := sc.LastPinged[conid]
		if pool, err := internal.GetStoragePool(conid); err == nil {
			if err = pool.Ping(); err != nil {
				sc.Close(conid, true)
				continue
			}
		}

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
		sc.Opened = lo.Filter[*containers.OpenedServerConnection](sc.Opened, func(x *containers.OpenedServerConnection, _ int) bool {
			return x.Conid != conid
		})
		sc.Closed[conid] = map[string]interface{}{
			"name":   "error",
			"status": existing.Status,
		}
		sc.LastPinged[conid] = 0
		_ = internal.DeleteStoragePool(conid)
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

func (sc *ServerConnections) pipeHandler(chData <-chan *containers.EchoMessage, conid string) {
	for {
		message, ok := <-chData
		if message != nil {
			if message.Err != nil {
				if existing := findByServerConnection(sc.Opened, conid); existing != nil && !existing.Disconnected {
					sc.Close(conid, true)
				}
			}
			switch message.MsgType {
			case "status":
				sc.handleStatus(conid, message.Payload.(*containers.OpenedStatus))
			case "version":
				sc.handleVersion(conid, message.Payload.(*standard.VersionMsg))
			case "databases":
				sc.handleDatabases(conid, message.Payload)
			case "ping":
				sc.handlePing()
			case "exit":
				break
			}
		}
		if !ok {
			break
		}
	}
}

func findByServerConnection(s []*containers.OpenedServerConnection, conid string) *containers.OpenedServerConnection {
	existing, ok := lo.Find[*containers.OpenedServerConnection](s, func(x *containers.OpenedServerConnection) bool {
		return x.Conid != "" && x.Conid == conid
	})

	if existing != nil && ok {
		return existing
	}
	return nil
}
