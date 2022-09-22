package sideQuests

import (
	"keeper/app/internal"
	"keeper/app/pkg/containers"
	"keeper/app/pkg/standard"
	"keeper/app/utility"
	"time"
)

var serverLastStatus string
var serverLastDatabases string
var serverLastPing utility.UnixTime

type StatusMessage struct {
	Name    string `json:"name"`
	Message string `json:"message"`
}

type ServerConnection struct {
}

func NewServerConnection() *ServerConnection {
	// setInterval(func() {
	// 	logger.Info("Server connection not alive, exiting")
	// 	ch <- &modules.EchoMessage{
	// 		MsgType: "exit",
	// 	}
	// })

	return &ServerConnection{}
}

/*
定时器 是当你想要在未来某一刻执行一次时使用的
打点器 则是当你想要在固定的时间间隔重复执行准备的。这里是一个打点器的例子，它将定时的执行，直到我们将它停止。
*/
func setInterval(fn func()) {
	ticker := time.NewTicker(time.Minute)
	go func(ticker *time.Ticker) {
		for range ticker.C {
			nowTime := utility.NewUnixTime()
			if nowTime-serverLastPing > utility.UnixTime(120*1000) {
				fn()
				ticker.Stop()
			}
		}
	}(ticker)
}

func setStatus(ch chan *containers.EchoMessage, data func() (*containers.OpenedStatus, error)) {
	status, err := data()
	if err != nil {
		ch <- &containers.EchoMessage{
			Err: err,
		}
		return
	}
	statusString := utility.ToJsonStr(status)
	if serverLastStatus != statusString {
		ch <- &containers.EchoMessage{Payload: status, MsgType: "status"}
		serverLastStatus = statusString
	}
}

func (msg *ServerConnection) ResetVars() {
	serverLastStatus = ""
	serverLastDatabases = ""
}

func (msg *ServerConnection) Connect(ch chan *containers.EchoMessage, conid string, connection map[string]interface{}) {
	defer close(ch)

	setStatus(ch, func() (*containers.OpenedStatus, error) {
		return &containers.OpenedStatus{Name: "pending"}, nil
	})

	driver, err := internal.TargetStoragePool(conid, connection)

	if err != nil {
		setStatus(ch, func() (*containers.OpenedStatus, error) {
			return &containers.OpenedStatus{Name: "error", Message: err.Error()}, err
		})
		return
	}

	if err = msg.readVersion(ch, driver); err != nil {
		return
	}

	if err = msg.handleRefresh(ch, driver); err != nil {
		return
	}

	//tasks.SetDriverPool(connection[conidkey].(string), sqlDriver)
	// ticker := time.NewTicker(2 * time.Second)
	// go func(ticker *time.Ticker) {
	// 	for range ticker.C {
	// 		if err := sqlDriver.Ping(); err != nil {
	// 			ch <- &modules.EchoMessage{
	// 				Payload: nil,
	// 				MsgType: "exit",
	// 				Dialect: sqlDriver.Dialect(),
	// 				Conid:   conid,
	// 			}
	// 			ticker.Stop()
	// 		}
	// 	}
	// }(ticker)
}

func (msg *ServerConnection) Ping() {
	serverLastPing = utility.NewUnixTime()
}

func (msg *ServerConnection) CreateDatabase() {

}

func (msg *ServerConnection) readVersion(ch chan *containers.EchoMessage, driver standard.SqlStandard) error {
	version, err := driver.GetVersion()
	if err != nil {
		setStatus(ch, func() (*containers.OpenedStatus, error) {
			return &containers.OpenedStatus{Name: "error", Message: err.Error()}, err
		})
		return err
	}

	ch <- &containers.EchoMessage{
		Payload: version,
		MsgType: "version",
	}

	return nil
}

func (msg *ServerConnection) handleRefresh(ch chan *containers.EchoMessage, driver standard.SqlStandard) error {
	databases, err := driver.ListDatabases()
	if err != nil {
		setStatus(ch, func() (*containers.OpenedStatus, error) {
			return &containers.OpenedStatus{Name: "error", Message: err.Error()}, err
		})
		return err
	}

	setStatus(ch, func() (*containers.OpenedStatus, error) {
		return &containers.OpenedStatus{Name: "ok"}, nil
	})

	databasesString := utility.ToJsonStr(databases)
	if serverLastDatabases != databasesString {
		ch <- &containers.EchoMessage{
			Payload: databases,
			MsgType: "databases",
			Dialect: driver.Dialect(),
		}
		serverLastDatabases = databasesString
	}
	return nil
}
