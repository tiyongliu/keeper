package sideQuests

import (
	"context"
	"errors"
	"github.com/samber/lo"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"keeper/app/db"
	"keeper/app/db/adapter"
	"keeper/app/db/adapter/mongo"
	"keeper/app/db/persist"
	"keeper/app/db/standard/modules"
	"keeper/app/internal/explorer"
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

func (msg *DatabaseConnection) Connect(ch chan *explorer.EchoMessage, newOpened *explorer.OpenedDatabaseConnection, structure interface{}) {
	defer close(ch)
	databaseLast = utility.NewUnixTime()
	if structure == nil {
		msg.setStatus(ch, func() (*explorer.OpenedStatus, error) {
			return &explorer.OpenedStatus{Name: "pending"}, nil
		})
	}
	driver, err := persist.GetStorageSession().Scanner(
		newOpened.Conid,
		lo.Assign(newOpened.Connection, map[string]interface{}{"database": newOpened.Database}),
	)

	if err != nil {
		msg.setStatus(ch, func() (*explorer.OpenedStatus, error) {
			return &explorer.OpenedStatus{Name: "error", Message: err.Error()}, err
		})
		return
	}

	version, err := readVersion(driver)

	if err != nil {
		msg.setStatus(ch, func() (*explorer.OpenedStatus, error) {
			return &explorer.OpenedStatus{Name: "error", Message: err.Error()}, err
		})
		return
	} else {
		ch <- &explorer.EchoMessage{
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

func (msg *DatabaseConnection) setStatus(ch chan *explorer.EchoMessage, data func() (*explorer.OpenedStatus, error)) {
	status, err := data()
	if err != nil {
		setStatus(ch, func() (*explorer.OpenedStatus, error) {
			return &explorer.OpenedStatus{Name: "error", Message: err.Error()}, nil
		})
		return
	}
	statusString := utility.ToJsonStr(status)
	if lastStatusString != statusString {
		status.Counter = getStatusCounter()
		ch <- &explorer.EchoMessage{MsgType: "status", Payload: status}
		lastStatusString = statusString
	}
}

func (msg *DatabaseConnection) readVersion(ch chan *explorer.EchoMessage, driver db.Session) error {
	version, err := driver.Version()
	if err != nil {
		setStatus(ch, func() (*explorer.OpenedStatus, error) {
			return &explorer.OpenedStatus{Name: "error", Message: err.Error()}, nil
		})
		return err
	}

	ch <- &explorer.EchoMessage{
		Payload: version,
		MsgType: "version",
		Dialect: driver.Dialect(),
	}

	return nil
}

func (msg *DatabaseConnection) handleFullRefresh(ch chan *explorer.EchoMessage, driver db.Session, database string) {
	loadingModel = true

	msg.setStatus(ch, func() (*explorer.OpenedStatus, error) {
		return &explorer.OpenedStatus{Name: "loadStructure"}, nil
	})

	analysedStructure = adapter.AnalyseFull(driver, database)
	analysedTime = utility.NewUnixTime()
	ch <- &explorer.EchoMessage{MsgType: "structure", Payload: analysedStructure}
	ch <- &explorer.EchoMessage{MsgType: "structureTime", Payload: analysedTime}
	msg.setStatus(ch, func() (*explorer.OpenedStatus, error) {
		return &explorer.OpenedStatus{Name: "ok"}, nil
	})

	loadingModel = false
}

func (msg *DatabaseConnection) handleIncrementalRefresh(ch chan *explorer.EchoMessage, forceSend bool, pool db.Session, database string) {
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
	ctx context.Context, conn *explorer.OpenedDatabaseConnection, selectParams interface{}) *explorer.EchoMessage {
	ch := make(chan *explorer.EchoMessage, 2)
	runtime.EventsEmit(ctx, "handleSqlSelect", selectParams)
	runtime.EventsOnce(ctx, "handleSqlSelectReturn", func(sql ...interface{}) {
		utility.WithRecover(func() {
			driver, err := persist.GetStorageSession().GetItem(conn.Conid, conn.Database)
			if err != nil {
				ch <- &explorer.EchoMessage{
					MsgType: "response",
					Err:     err,
				}
				return
			}
			ch <- msg.handleQueryData(driver, sql[0].(string), true)
		}, func(err error) {
			ch <- &explorer.EchoMessage{
				MsgType: "response",
				Err:     errors.New(serializer.ErrNil),
			}
		})
	})
	defer close(ch)
	return <-ch
}

func (msg *DatabaseConnection) handleQueryData(driver db.Session, sql string, skipReadonlyCheck bool) *explorer.EchoMessage {
	res, err := driver.Query(sql)
	return &explorer.EchoMessage{
		Payload: res,
		MsgType: "response",
		Err:     err,
	}
}

func (msg *DatabaseConnection) HandleCollectionData(conn *explorer.OpenedDatabaseConnection,
	options *modules.CollectionDataOptions) *explorer.EchoMessage {
	driver, err := persist.GetStorageSession().GetItem(conn.Conid, conn.Database)
	if err != nil {
		return &explorer.EchoMessage{
			MsgType: "response",
			Err:     err,
		}
	}

	collection, err := driver.(*mongo.Source).ReadCollection(conn.Database, options)
	if err != nil {
		return &explorer.EchoMessage{
			MsgType: "response",
			Err:     err,
		}
	}

	if options != nil && options.CountDocuments {
		return &explorer.EchoMessage{
			Payload: map[string]interface{}{"count": collection},
			MsgType: "response",
		}
	} else {
		return &explorer.EchoMessage{
			Payload: map[string]interface{}{"rows": collection},
			MsgType: "response",
		}
	}
}

func (msg *DatabaseConnection) ReadVersion(ch chan *explorer.EchoMessage, driver db.Session) error {
	version, err := driver.Version()
	if err != nil {
		return err
	}

	ch <- &explorer.EchoMessage{
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

func handleDriverDataCore() {

}
