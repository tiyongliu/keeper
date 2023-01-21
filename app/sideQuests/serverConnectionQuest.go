package sideQuests

import (
	"keeper/app/db"
	"keeper/app/db/drivers"
	"keeper/app/internal/explorer"
	"keeper/app/utility"
	"time"
)

var serverLastStatus string
var serverLastDatabases string
var serverLastPing utility.UnixTime
var storedConnection = make(map[string]interface{})

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

func setStatus(ch chan *explorer.EchoMessage, data func() (*explorer.OpenedStatus, error)) {
	status, err := data()
	if err != nil {
		ch <- &explorer.EchoMessage{
			Err: err,
		}
		return
	}
	statusString := utility.ToJsonStr(status)
	if serverLastStatus != statusString {
		ch <- &explorer.EchoMessage{Payload: status, MsgType: "status"}
		serverLastStatus = statusString
	}
}

func (msg *ServerConnection) ResetVars() {
	serverLastStatus = ""
	serverLastDatabases = ""
}

func (msg *ServerConnection) Connect(ch chan *explorer.EchoMessage, connection map[string]interface{}) {
	storedConnection = connection
	defer close(ch)

	setStatus(ch, func() (*explorer.OpenedStatus, error) {
		return &explorer.OpenedStatus{Name: "pending"}, nil
	})

	driver, err := drivers.NewCompatDriver().Open(connection)
	if err != nil {
		setStatus(ch, func() (*explorer.OpenedStatus, error) {
			return &explorer.OpenedStatus{Name: "error", Message: err.Error()}, err
		})
		return
	}
	//targetDriver, err := internal.RequireEngineDriver(connection)
	//if err != nil {
	//
	//}
	//internal.ConnectUtility1(targetDriver, storedConnection, "app")

	//driver, err := internal.TargetStoragePool(conid, connection)
	//connectUtility

	if err != nil {
		setStatus(ch, func() (*explorer.OpenedStatus, error) {
			return &explorer.OpenedStatus{Name: "error", Message: err.Error()}, err
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

func (msg *ServerConnection) readVersion(ch chan *explorer.EchoMessage, driver db.Session) error {
	version, err := driver.Version()
	if err != nil {
		setStatus(ch, func() (*explorer.OpenedStatus, error) {
			return &explorer.OpenedStatus{Name: "error", Message: err.Error()}, err
		})
		return err
	}

	ch <- &explorer.EchoMessage{
		Payload: version,
		MsgType: "version",
	}

	return nil
}

func (msg *ServerConnection) handleRefresh(ch chan *explorer.EchoMessage, driver db.Session) error {
	databases, err := driver.ListDatabases()
	if err != nil {
		setStatus(ch, func() (*explorer.OpenedStatus, error) {
			return &explorer.OpenedStatus{Name: "error", Message: err.Error()}, err
		})
		return err
	}

	setStatus(ch, func() (*explorer.OpenedStatus, error) {
		return &explorer.OpenedStatus{Name: "ok"}, nil
	})

	databasesString := utility.ToJsonStr(databases)
	if serverLastDatabases != databasesString {
		ch <- &explorer.EchoMessage{
			Payload: databases,
			MsgType: "databases",
			Dialect: driver.Dialect(),
		}
		serverLastDatabases = databasesString
	}
	return nil
}
