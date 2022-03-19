package sideQuests

import (
<<<<<<< HEAD
	"keeper/app/internal"
	"keeper/app/pkg/containers"
=======
	"keeper/app/code"
	"keeper/app/modules"
>>>>>>> 90ec4d6 (数据库连接)
	"keeper/app/pkg/standard"
	"keeper/app/utility"
	"time"
)

var ServerLastStatus string
var serverLastPing utility.UnixTime

var ServerLastDatabases string

type StatusMessage struct {
	Name    string `json:"name"`
	Message string `json:"message"`
}

type ServerConnection struct {
}

<<<<<<< HEAD
func NewServerConnection() *ServerConnection {
=======
func NewServerConnection(ch chan *modules.EchoMessage) *ServerConnection {
>>>>>>> 90ec4d6 (数据库连接)
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

<<<<<<< HEAD
func setStatus(ch chan *containers.EchoMessage, data func() (*containers.OpenedStatus, error)) {
	status, err := data()
	if err != nil {
<<<<<<< HEAD
		ch <- &containers.EchoMessage{
			Err: err,
=======
		return
	}

	//TODO connectUtility, 可以传递一个func 因为返回值都是一样的，在func内部进行处理
	var driver standard.SqlStandard
=======
func (msg *ServerConnection) Connect(ch chan *modules.EchoMessage, conid string, connection map[string]interface{}) {
	msg.setStatus(ch, conid, "pending")
	serverLastPing = tools.NewUnixTime()
	//TODO connectUtility, 可以传递一个func 因为返回值都是一样的，在func内部进行处理
	sqlDriver, err := GetSqlDriver(connection)
	if err != nil {
		msg.setStatus(ch, conid, "error", err.Error())
		return
	}

	msg.SqlDriver = sqlDriver

	if err := msg.readVersion(ch, conid, sqlDriver); err != nil {
		msg.setStatus(ch, conid, "error", err.Error())
		return
	}

	if err := msg.handleRefresh(ch, conid, sqlDriver); err != nil {
		msg.setStatus(ch, conid, "error", err.Error())
		return
	}

	msg.setStatus(ch, conid, "ok")

	ch <- &modules.EchoMessage{
		MsgType: "pool",
		Payload: sqlDriver,
		Conid:   conid,
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

func GetSqlDriver(connection map[string]interface{}) (driver standard.SqlStandard, err error) {
>>>>>>> 90ec4d6 (数据库连接)
	switch connection["engine"].(string) {
	case code.MYSQLALIAS:
		driver, err = NewMysqlDriver(connection)
		if err != nil {
			return
		}
		msg.Mysql = driver
	case code.MONGOALIAS:
		driver, err = NewMongoDriver(connection)
		if err != nil {
			return
>>>>>>> 35a8c6f (临时编写了readme)
		}
		return
	}
	statusString := utility.ToJsonStr(status)
	if ServerLastStatus != statusString {
		ch <- &containers.EchoMessage{Payload: status, MsgType: "status"}
		ServerLastStatus = statusString
	}
}

<<<<<<< HEAD
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
=======
func (msg *ServerConnection) setStatus(ch chan *modules.EchoMessage, conid, name string, message ...string) {
	status := map[string]string{"name": name}
	if len(message) > 0 {
		status["message"] = message[0]
	}

	statusString := tools.ToJsonStr(status)
	if serverLastStatus != statusString {
		ch <- &modules.EchoMessage{
			MsgType: "status",
			Payload: status,
			Conid:   conid,
		}
		serverLastStatus = statusString
	}

	if name == "error" {
		ch <- &modules.EchoMessage{
			MsgType: "exit",
			Conid:   conid,
		}
>>>>>>> 90ec4d6 (数据库连接)
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

<<<<<<< HEAD
func (msg *ServerConnection) readVersion(ch chan *containers.EchoMessage, driver standard.SqlStandard) error {
	version, err := driver.GetVersion()
=======
func (msg *ServerConnection) readVersion(ch chan *modules.EchoMessage, conid string, pool standard.SqlStandard) error {
	version, err := pool.GetVersion()
>>>>>>> 90ec4d6 (数据库连接)
	if err != nil {
		setStatus(ch, func() (*containers.OpenedStatus, error) {
			return &containers.OpenedStatus{Name: "error", Message: err.Error()}, err
		})
		return err
	}

	ch <- &containers.EchoMessage{
		Payload: version,
		MsgType: "version",
		Conid:   conid,
	}

	return nil
}

<<<<<<< HEAD
func (msg *ServerConnection) handleRefresh(ch chan *containers.EchoMessage, driver standard.SqlStandard) error {
	databases, err := driver.ListDatabases()
=======
func (msg *ServerConnection) handleRefresh(ch chan *modules.EchoMessage, conid string, pool standard.SqlStandard) error {
	databases, err := pool.ListDatabases()
	msg.setStatus(ch, conid, "ok")
	databasesString := tools.ToJsonStr(databases)
>>>>>>> 90ec4d6 (数据库连接)
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
	if ServerLastDatabases != databasesString {
		ch <- &containers.EchoMessage{
			Payload: databases,
			MsgType: "databases",
<<<<<<< HEAD
			Dialect: driver.Dialect(),
=======
			Dialect: pool.Dialect(),
			Conid:   conid,
>>>>>>> 90ec4d6 (数据库连接)
		}
		ServerLastDatabases = databasesString
	}
	ch <- &containers.EchoMessage{
		Payload: nil,
		MsgType: "exit",
		Dialect: driver.Dialect(),
	}
	return nil
}
