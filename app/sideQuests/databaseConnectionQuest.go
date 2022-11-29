package sideQuests

import (
	"context"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"keeper/app/adapter"
	"keeper/app/internal"
	"keeper/app/pkg/containers"
	"keeper/app/pkg/logger"
	"keeper/app/pkg/standard"
	"keeper/app/utility"
)

var databaseLast utility.UnixTime

var analysedStructure map[string]interface{}

var analysedTime utility.UnixTime = 0
var loadingModel bool
var statusCounter int
var lastStatusString string

func getStatusCounter() int {
	statusCounter += 1
	return statusCounter
}

type DatabaseConnection struct {
}

func NewDatabaseConnection() *DatabaseConnection {
	return &DatabaseConnection{}
}

func (msg *DatabaseConnection) ResetVars() {
	lastStatusString = ""
}

func (msg *DatabaseConnection) Connect(ch chan *containers.EchoMessage, newOpened *containers.OpenedDatabaseConnection, structure interface{}) {
	defer close(ch)
	databaseLast = utility.NewUnixTime()
	if structure == nil {
		msg.setStatus(ch, func() (*containers.OpenedStatus, error) {
			return &containers.OpenedStatus{Name: "pending"}, nil
		})
	}

	driver, err := internal.TargetStoragePool(newOpened.Conid, newOpened.Connection)
	if err != nil {
		msg.setStatus(ch, func() (*containers.OpenedStatus, error) {
			return &containers.OpenedStatus{Name: "error", Message: err.Error()}, err
		})
		return
	}

	version, err := readVersion(driver)

	if err != nil {
		msg.setStatus(ch, func() (*containers.OpenedStatus, error) {
			return &containers.OpenedStatus{Name: "error", Message: err.Error()}, err
		})
		return
	} else {
		ch <- &containers.EchoMessage{
			Payload: version,
			MsgType: "version",
		}
	}

	if structure != nil {
		msg.handleIncrementalRefresh(ch, true, driver, newOpened.Database)
	} else {
	}

	msg.handleFullRefresh(ch, driver, newOpened.Database)
}

func (msg *DatabaseConnection) setStatus(ch chan *containers.EchoMessage, data func() (*containers.OpenedStatus, error)) {
	status, err := data()
	if err != nil {
		ch <- &containers.EchoMessage{Err: err}
		return
	}
	statusString := utility.ToJsonStr(status)
	if lastStatusString != statusString {
		status.Counter = getStatusCounter()
		ch <- &containers.EchoMessage{MsgType: "status", Payload: status}
		lastStatusString = statusString
	}
}

func (msg *DatabaseConnection) readVersion(ch chan *containers.EchoMessage, pool standard.SqlStandard) error {
	version, err := pool.GetVersion()
	if err != nil {
		return err
	}

	ch <- &containers.EchoMessage{
		Payload: version,
		MsgType: "version",
		Dialect: pool.Dialect(),
	}

	return nil
}

func (msg *DatabaseConnection) handleFullRefresh(ch chan *containers.EchoMessage, pool standard.SqlStandard, database string) {
	loadingModel = true

	msg.setStatus(ch, func() (*containers.OpenedStatus, error) {
		return &containers.OpenedStatus{Name: "loadStructure"}, nil
	})

	analysedStructure = adapter.AnalyseFull(pool, database)
	analysedTime = utility.NewUnixTime()
	ch <- &containers.EchoMessage{MsgType: "structure", Payload: analysedStructure}
	ch <- &containers.EchoMessage{MsgType: "structureTime", Payload: analysedTime}
	msg.setStatus(ch, func() (*containers.OpenedStatus, error) {
		return &containers.OpenedStatus{Name: "ok"}, nil
	})

	loadingModel = false
}

func (msg *DatabaseConnection) handleIncrementalRefresh(ch chan *containers.EchoMessage, forceSend bool, pool standard.SqlStandard, database string) {
	/*loadingModel = true
	msg.setStatus(ch, func() (*containers.OpenedStatus, error) {
		return &containers.OpenedStatus{Name: "checkStructure", Counter: getStatusCounter()}, nil
	})

	tables, err := pool.Tables(args...)
	if err != nil {
		setStatus(ch, func() (*containers.OpenedStatus, error) {
			return &containers.OpenedStatus{Name: "error", Message: err.Error()}, err
		})
		return
	}
	analysedTime = utility.NewUnixTime()

	if forceSend || tables != nil {
		ch <- &containers.EchoMessage{
			MsgType: "structure",
			Payload: map[string]interface{}{
				"collections": tables,
				"engine":      pool.Dialect(),
			},
		}
	}

	ch <- &containers.EchoMessage{
		MsgType: "structureTime",
		Payload: analysedTime,
	}

	msg.setStatus(ch, func() (*containers.OpenedStatus, error) {
		return &containers.OpenedStatus{Name: "ok"}, nil
	})*/
}

func (msg *DatabaseConnection) HandleSqlSelect(ch chan *containers.EchoMessage, ctx context.Context, conn *containers.OpenedDatabaseConnection, msgid string, _select map[string]interface{}) (err error) {
	runtime.EventsEmit(ctx, "handleSqlSelect", _select)
	runtime.EventsOn(ctx, "handleSqlSelectReturn", func(sql ...interface{}) {
		utility.WithRecover(func() {
			driver, e := internal.GetStoragePool(conn.Conid)
			if err != nil {
				err = e
			}
			logger.Info("11111111111")
			msg.handleQueryData(driver, msgid, sql[0].(string), true)
		}, func(er error) {
			err = er
		})
	})
	logger.Info("22222222222")
	return err
}

func (msg *DatabaseConnection) handleQueryData(driver standard.SqlStandard, msgid, sql string, skipReadonlyCheck bool) {
	//driver
	//get driver
	//driver query

	driver.Query(sql)
}

func (msg *DatabaseConnection) ReadVersion(ch chan *containers.EchoMessage, pool standard.SqlStandard) error {
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

func (msg *DatabaseConnection) QueryData() {

}

func (msg *DatabaseConnection) SyncModel() {
	if loadingModel {
		return
	}
}

func (msg *DatabaseConnection) Ping() {
	databaseLast = utility.NewUnixTime()
}

func HandleSqlSelect() {

}
