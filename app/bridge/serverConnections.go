package bridge

import (
	"fmt"
<<<<<<< HEAD
	"github.com/samber/lo"
	"keeper/app/internal"
	"keeper/app/pkg/containers"
=======
	"keeper/app/code"
	"keeper/app/modules"
	"keeper/app/pkg/logger"
>>>>>>> 90ec4d6 (数据库连接)
	"keeper/app/pkg/serializer"
	"keeper/app/pkg/standard"
	"keeper/app/sideQuests"
	"keeper/app/utility"
	"sync"
	"time"

	"github.com/samber/lo"
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
<<<<<<< HEAD
=======
		// PoolMap:                 make(map[string]standard.SqlStandard),
>>>>>>> 90ec4d6 (数据库连接)
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

<<<<<<< HEAD
	existing.Databases = databases
=======
	existing["databases"] = databases

>>>>>>> 90ec4d6 (数据库连接)
	utility.EmitChanged(Application.ctx, fmt.Sprintf("database-list-changed-%s", conid))
}

func (sc *ServerConnections) handleVersion(conid string, version *standard.VersionMsg) {
<<<<<<< HEAD
	existing := findByServerConnection(sc.Opened, conid)

	if existing == nil {
		return
	}

	existing.Version = version
=======
	existing, ok := lo.Find[map[string]interface{}](sc.Opened, func(x map[string]interface{}) bool {
		uuid, ok := x[conidkey].(string)
		return ok && uuid == conid
	})

	if existing == nil || !ok {
		return
	}

	existing["version"] = version
>>>>>>> 90ec4d6 (数据库连接)
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

<<<<<<< HEAD
	if existing != nil {
		utility.EmitChanged(Application.ctx, "server-status-changed")
=======
	existing, ok := lo.Find[map[string]interface{}](sc.Opened, func(x map[string]interface{}) bool {
		uuid, ok := x[conidkey].(string)
		return ok && uuid == conid
	})

	utility.EmitChanged(Application.ctx, "server-status-changed")
	if existing != nil && ok {
		if err := sc.checker(conid); err != nil {
			existing["status"] = &modules.OpenedStatus{Name: "error", Message: err.Error()}
			sc.Close(conid, true)
		}
>>>>>>> 90ec4d6 (数据库连接)
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

<<<<<<< HEAD
	ch := make(chan *containers.EchoMessage)
	utility.EmitChanged(Application.ctx, "server-status-changed")

	defer func() {
		go sc.ServerConnectionChannel.Connect(ch, conid, connection)
		go sc.pipeHandler(ch, conid)
	}()
=======
	go sc.ServerConnectionChannel.Connect(sc.ch, conid, connection)
	go sc.pipeHandler(newOpened, sc.ch)
>>>>>>> 90ec4d6 (数据库连接)

	return newOpened
}

func (sc *ServerConnections) checker(conid string) error {
	pool, err := internal.GetDriverPool(conid)
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
	time.Sleep(time.Millisecond * 100)
	values := map[string]interface{}{}
<<<<<<< HEAD
	for _, driver := range sc.Opened {
		values[driver.Conid] = driver.Status
	}
	for key, val := range sc.Closed {
		values[key] = val
	}
=======

	for _, driver := range sc.Opened {
		values[driver[conidkey].(string)] = driver["status"]
	}

	for key, val := range sc.Closed {
		values[key] = val
	}

>>>>>>> 90ec4d6 (数据库连接)
	return serializer.SuccessData("", values)
}

func (sc *ServerConnections) Ping(connections []string) *serializer.Response {
	for _, conid := range lo.Uniq[string](connections) {
		last := sc.LastPinged[conid]
<<<<<<< HEAD
		if pool, err := internal.GetDriverPool(conid); err == nil {
			if err = pool.Ping(); err != nil {
				sc.Close(conid, true)
				continue
			}
=======
		if last > 0 && tools.NewUnixTime()-last < tools.GetUnixTime(30*1000) {
			continue
>>>>>>> 90ec4d6 (数据库连接)
		}

		if last > 0 && utility.NewUnixTime()-last < utility.GetUnixTime(30*1000) {
			continue
		}

		sc.LastPinged[conid] = utility.NewUnixTime()
		sc.ensureOpened(conid)
		sc.ServerConnectionChannel.Ping()
	}

	return serializer.SuccessData("", map[string]string{"status": "ok"})
}

func (sc *ServerConnections) Close(conid string, kill bool) {
<<<<<<< HEAD
	existing := findByServerConnection(sc.Opened, conid)
	if existing != nil {
		existing.Disconnected = true
=======
	existing, ok := lo.Find[map[string]interface{}](sc.Opened, func(item map[string]interface{}) bool {
		return item[conidkey].(string) != "" && item[conidkey].(string) == conid
	})

	if existing != nil && ok {
		existing["disconnected"] = true
>>>>>>> 90ec4d6 (数据库连接)
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

		internal.DeleteDriverPool(conid)

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

func (sc *ServerConnections) pipeHandler(chData <-chan *containers.EchoMessage, conid string) {
	for {
		message, ok := <-chData
<<<<<<< HEAD
=======
		logger.Infof("current: %s", message.MsgType)
		conid := message.Conid
>>>>>>> 90ec4d6 (数据库连接)
		if message != nil {
			if message.Err != nil {
				if existing := findByServerConnection(sc.Opened, conid); existing != nil && !existing.Disconnected {
					sc.Close(conid, true)
				}
			}
			switch message.MsgType {
			case "status":
<<<<<<< HEAD
				sc.handleStatus(conid, message.Payload.(*containers.OpenedStatus))
=======
				sc.handleStatus(conid, message.Payload.(map[string]string))
>>>>>>> 90ec4d6 (数据库连接)
			case "version":
				sc.handleVersion(conid, message.Payload.(*standard.VersionMsg))
			case "databases":
				sc.handleDatabases(conid, message.Payload)
			case "ping":
				sc.handlePing()
			case "exit":
<<<<<<< HEAD
=======
				if !newOpened["disconnected"].(bool) {
					sc.Close(conid, true)
				}
			case "pool":
				if sc.PoolMap == nil {
					sc.PoolMap = map[string]standard.SqlStandard{conid: message.Payload.(standard.SqlStandard)}
				} else {
					sc.PoolMap[conid] = message.Payload.(standard.SqlStandard)
				}
			}

			if !ok {
				if !newOpened["disconnected"].(bool) {
					sc.Close(conid, true)
				}
>>>>>>> 90ec4d6 (数据库连接)
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
