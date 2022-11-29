package sideQuests

import (
	"context"
	"errors"
	"github.com/samber/lo"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"keeper/app/db"
	"keeper/app/db/adapter"
	"keeper/app/db/persist"
	"keeper/app/pkg/containers"
	"keeper/app/pkg/serializer"
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

	//logger.Infof("newOpened Connect req: %s", utility.ToJsonStr(lo.Assign(newOpened.Connection, map[string]interface{}{"database": newOpened.Database})))
	driver, err := persist.GetStorageSession().Scanner(
		newOpened.Conid,
		lo.Assign(newOpened.Connection, map[string]interface{}{"database": newOpened.Database}),
	)
	//driver, err := internal.TargetStoragePool(newOpened.Conid, newOpened.Connection)

	//driver, err := drivers.NewCompatDriver().Open(
	//	lo.Assign(newOpened.Connection, map[string]interface{}{"database": newOpened.Database}),
	//)

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

func (msg *DatabaseConnection) readVersion(ch chan *containers.EchoMessage, driver db.Session) error {
	version, err := driver.Version()
	if err != nil {
		return err
	}

	ch <- &containers.EchoMessage{
		Payload: version,
		MsgType: "version",
		Dialect: driver.Dialect(),
	}

	return nil
}

func (msg *DatabaseConnection) handleFullRefresh(ch chan *containers.EchoMessage, driver db.Session, database string) {
	loadingModel = true

	msg.setStatus(ch, func() (*containers.OpenedStatus, error) {
		return &containers.OpenedStatus{Name: "loadStructure"}, nil
	})

	analysedStructure = adapter.AnalyseFull(driver, database)
	analysedTime = utility.NewUnixTime()
	ch <- &containers.EchoMessage{MsgType: "structure", Payload: analysedStructure}
	ch <- &containers.EchoMessage{MsgType: "structureTime", Payload: analysedTime}
	msg.setStatus(ch, func() (*containers.OpenedStatus, error) {
		return &containers.OpenedStatus{Name: "ok"}, nil
	})

	loadingModel = false
}

func (msg *DatabaseConnection) handleIncrementalRefresh(ch chan *containers.EchoMessage, forceSend bool, pool db.Session, database string) {
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

func (msg *DatabaseConnection) HandleSqlSelect(
	ctx context.Context, conn *containers.OpenedDatabaseConnection, selectParams interface{}) *containers.EchoMessage {
	ch := make(chan *containers.EchoMessage, 1)
	runtime.EventsEmit(ctx, "handleSqlSelect", selectParams)
	runtime.EventsOnce(ctx, "handleSqlSelectReturn", func(sql ...interface{}) {
		utility.WithRecover(func() {
			driver, err := persist.GetStorageSession().GetItem(conn.Conid, conn.Database)
			if err != nil {
				ch <- &containers.EchoMessage{
					MsgType: "response",
					Err:     err,
				}
				return
			}

			//todo 这个短期测试使用，后期需要删除掉
			if driver.Ping() != nil {
				ch <- &containers.EchoMessage{
					MsgType: "response",
					Err:     driver.Ping(),
				}
				return
			}
			ch <- msg.handleQueryData(driver, sql[0].(string), true)
		}, func(err error) {
			ch <- &containers.EchoMessage{
				MsgType: "response",
				Err:     errors.New(serializer.ErrNil),
			}
		})
	})
	defer close(ch)
	return <-ch
}

func (msg *DatabaseConnection) handleQueryData(driver db.Session, sql string, skipReadonlyCheck bool) *containers.EchoMessage {
	res, err := driver.Query(sql)
	return &containers.EchoMessage{
		Payload: res,
		MsgType: "response",
		Err:     err,
	}
}

func (msg *DatabaseConnection) ReadVersion(ch chan *containers.EchoMessage, driver db.Session) error {
	version, err := driver.Version()
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
