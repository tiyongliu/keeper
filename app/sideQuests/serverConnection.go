package sideQuests

import (
	"keeper/app/pkg/containers"
	"keeper/app/pkg/standard"
	"keeper/app/utility"
	"time"
)

var ServerLastStatus string
var serverLastPing utility.UnixTime

var ServerlastDatabases string

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
		close(ch)
		return
	}
	statusString := utility.ToJsonStr(status)
	if ServerLastStatus != statusString {
		ch <- &containers.EchoMessage{Payload: status, MsgType: "status"}
		ServerLastStatus = statusString
	}
}

func (msg *ServerConnection) Connect(ch chan *containers.EchoMessage, connectUtility func() (driver standard.SqlStandard, err error)) {
	setStatus(ch, func() (*containers.OpenedStatus, error) {
		return &containers.OpenedStatus{Name: "pending"}, nil
	})

	sqlDriver, err := connectUtility()
	if err != nil {
		setStatus(ch, func() (*containers.OpenedStatus, error) {
			return &containers.OpenedStatus{Name: "error", Message: err.Error()}, err
		})
		return
	}

	if err := msg.readVersion(ch, sqlDriver); err != nil {
		setStatus(ch, func() (*containers.OpenedStatus, error) {
			return &containers.OpenedStatus{Name: "error", Message: err.Error()}, err
		})
		return
	}

	if err := msg.handleRefresh(ch, sqlDriver); err != nil {
		return
	}

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

//func (msg *ServerConnection) Connect(ch chan *modules.EchoMessage, connection map[string]interface{}) {
//	setStatus(ch, func() (*modules.OpenedStatus, error) {
//		return &modules.OpenedStatus{Name: "ping"}, nil
//	})
//	serverLastPing = tools.NewUnixTime()
//	//TODO connectUtility, 可以传递一个func 因为返回值都是一样的，在func内部进行处理
//	sqlDriver, err := GetSqlDriver(connection)
//	if err != nil {
//		setStatus(ch, func() (*modules.OpenedStatus, error) {
//			return &modules.OpenedStatus{Name: "error", Message: err.Error()}, err
//		})
//		return
//	}
//
//	if err := msg.readVersion(ch, sqlDriver); err != nil {
//		setStatus(ch, func() (*modules.OpenedStatus, error) {
//			return &modules.OpenedStatus{Name: "error", Message: err.Error()}, err
//		})
//		return
//	}
//
//	if err := msg.handleRefresh(ch, sqlDriver); err != nil {
//		return
//	}
//
//}

func (msg *ServerConnection) NewTime() {
	serverLastPing = utility.NewUnixTime()
}

func (msg *ServerConnection) CreateDatabase() {

}

func SendChanMessage(ch chan *containers.EchoMessage, data func() (*containers.EchoMessage, error)) {
	message, err := data()
	if err != nil {
		close(ch)
		return
	}
	ch <- message
}

func (msg *ServerConnection) setStatus(ch chan *containers.EchoMessage, name string, message ...string) {
	status := &containers.OpenedStatus{Name: name}
	if len(message) > 0 {
		status.Message = message[0]
	}

	statusString := utility.ToJsonStr(status)
	if ServerLastStatus != statusString {
		ch <- &containers.EchoMessage{
			MsgType: "status",
			Payload: status,
		}
		ServerLastStatus = statusString
	}

	if name == "error" {
		ch <- &containers.EchoMessage{
			MsgType: "exit",
		}
	}
}

func (msg *ServerConnection) readVersion(ch chan *containers.EchoMessage, pool standard.SqlStandard) error {
	version, err := pool.GetVersion()
	if err != nil {
		return err
	}

	ch <- &containers.EchoMessage{
		Payload: version,
		MsgType: "version",
	}

	return nil
}

func (msg *ServerConnection) handleRefresh(ch chan *containers.EchoMessage, pool standard.SqlStandard) error {
	databases, err := pool.ListDatabases()
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
	if ServerlastDatabases != databasesString {
		ch <- &containers.EchoMessage{
			Payload: databases,
			MsgType: "databases",
			Dialect: pool.Dialect(),
		}
		ServerlastDatabases = databasesString
	}
	ch <- &containers.EchoMessage{
		Payload: nil,
		MsgType: "exit",
		Dialect: pool.Dialect(),
	}
	return nil
}
