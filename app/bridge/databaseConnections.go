package bridge

import (
	"fmt"
	"github.com/samber/lo"
	uuid "github.com/satori/go.uuid"
	"keeper/app/pkg/containers"
	"keeper/app/pkg/logger"
	"keeper/app/pkg/serializer"
	"keeper/app/pkg/standard"
	"keeper/app/plugins"
	"keeper/app/sideQuests"
	"keeper/app/utility"
)

const databaseKey = "database"

type DatabaseConnections struct {
	Opened             []*containers.OpenedDatabaseConnection
	Closed             map[string]*containers.DatabaseConnectionClosed
	DatabaseConnection *sideQuests.DatabaseConnection
}

func NewDatabaseConnections() *DatabaseConnections {
	return &DatabaseConnections{
		Closed:             make(map[string]*containers.DatabaseConnectionClosed),
		DatabaseConnection: sideQuests.NewDatabaseConnection(),
	}
}

func (dc *DatabaseConnections) handleStructure(conid, database string, structure map[string]interface{}) {
	existing := findByDatabaseConnection(dc.Opened, conid, database)

	if existing == nil {
		return
	}

	existing.Structure = structure
	utility.EmitChanged(Application.ctx, fmt.Sprintf("database-structure-changed-%s-%s", conid, database))
}

func (dc *DatabaseConnections) handleStructureTime(conid, database string, analysedTime utility.UnixTime) {
	existing := findByDatabaseConnection(dc.Opened, conid, database)

	if existing == nil {
		return
	}

	existing.AnalysedTime = analysedTime

	utility.EmitChanged(Application.ctx, fmt.Sprintf("database-status-changed-%s-%s", conid, database))
}

func (dc *DatabaseConnections) handleVersion(conid, database string, version *standard.VersionMsg) {
	existing := findByDatabaseConnection(dc.Opened, conid, database)
	if existing == nil {
		return
	}
	existing.ServerVersion = version
	utility.EmitChanged(Application.ctx, fmt.Sprintf("database-server-version-changed-%s-%s", conid, database))
}

func (dc *DatabaseConnections) handleError(conid, database string, err error) {
	logger.Errorf("Error in database connection [%s], database [%d]: [%v]", conid, database, err)
}

func (dc *DatabaseConnections) handleStatus(conid, database string, status *containers.OpenedStatus) {
	existing := findByDatabaseConnection(dc.Opened, conid, database)
	if existing == nil {
		return
	}
	if existing.Status != nil && status != nil && existing.Status.Counter > status.Counter {
		//logger.Infof("456--- %d", existing.Status.Counter)
		return
	}
	existing.Status = status
	utility.EmitChanged(Application.ctx, fmt.Sprintf("database-status-changed-%s-%s", conid, database))
}

func (dc *DatabaseConnections) handlePing() {

}

func (dc *DatabaseConnections) ensureOpened(conid, database string) *containers.OpenedDatabaseConnection {
	existing := findByDatabaseConnection(dc.Opened, conid, database)

	if existing != nil {
		return existing
	}

	connection := getCore(conid, false)
	lastClosed := dc.Closed[fmt.Sprintf("%s/%s", conid, database)]

	newOpened := &containers.OpenedDatabaseConnection{
		Conid:         conid,
		Status:        &containers.OpenedStatus{Name: "pending"},
		Database:      database,
		Connection:    connection,
		ServerVersion: nil,
	}

	if lastClosed == nil || lastClosed.Structure == nil {
		newOpened.Structure = plugins.CreateEmptyStructure()
	}

	dc.Opened = append(dc.Opened, newOpened)

	var structure map[string]interface{}
	if lastClosed != nil && lastClosed.Structure == nil {
		structure = lastClosed.Structure
	} else {
		structure = nil
	}

	ch := make(chan *containers.EchoMessage)
	defer func() {
		dc.DatabaseConnection.ResetVars()
		go dc.DatabaseConnection.Connect(ch, newOpened, structure)
		go dc.pipeHandler(ch, conid, database)
	}()

	return newOpened
}

func (dc *DatabaseConnections) Refresh(req *DatabaseKeepOpenRequest) *serializer.Response {
	if !req.KeepOpen {
		dc.close(req.Conid, req.Database, true)
	}
	dc.ensureOpened(req.Conid, req.Database)
	return serializer.SuccessData(serializer.SUCCESS, map[string]string{"status": "ok"})
}

func (dc *DatabaseConnections) SyncModel(req *DatabaseRequest) *serializer.Response {
	dc.ensureOpened(req.Conid, req.Database)

	return serializer.SuccessData(serializer.SUCCESS, map[string]string{"status": "ok"})
}

type databaseConnections struct {
	Conid    string `json:"conid"`
	Database string `json:"database"`
}

type DatabaseRequest struct {
	databaseConnections
}

type DatabaseKeepOpenRequest struct {
	databaseConnections
	KeepOpen bool `json:"keepOpen"`
}

func (dc *DatabaseConnections) Ping(req *DatabaseRequest) *serializer.Response {
	if req == nil || req.Conid == "" {
		return serializer.Fail(serializer.IdNotEmpty)
	}

	existing := findByDatabaseConnection(dc.Opened, req.Conid, req.Database)

	if existing != nil {
		dc.DatabaseConnection.Ping()
	} else {
		existing = dc.ensureOpened(req.Conid, req.Database)
	}

	res := map[string]interface{}{"status": "ok"}
	if existing != nil {
		res["connectionStatus"] = existing.Status
	} else {
		res["connectionStatus"] = nil
	}
	return serializer.SuccessData(serializer.SUCCESS, res)
}

func (dc *DatabaseConnections) Structure(req *DatabaseRequest) *serializer.Response {
	if req.Conid == "__model" {
		//todo  const model = await importDbModel(database);
	}

	opened := dc.ensureOpened(req.Conid, req.Database)
	return serializer.SuccessData(serializer.SUCCESS, opened.Structure)
}

func (dc *DatabaseConnections) pipeHandler(chData <-chan *containers.EchoMessage, conid, database string) {
	for {
		message, ok := <-chData
		if message != nil {
			if message.Err != nil {
				if existing := findByDatabaseConnection(dc.Opened, conid, database); existing != nil && !existing.Disconnected {
					dc.close(conid, database, false)
				}
			}
			switch message.MsgType {
			case "status":
				dc.handleStatus(conid, database, message.Payload.(*containers.OpenedStatus))
			case "structure":
				dc.handleStructure(conid, database, message.Payload.(map[string]interface{}))
			case "structureTime":
				dc.handleStructureTime(conid, database, message.Payload.(utility.UnixTime))
			case "version":
				dc.handleVersion(conid, database, message.Payload.(*standard.VersionMsg))
			}
		}
		if !ok {
			break
		}
	}
}

func (dc *DatabaseConnections) ServerVersion(req *DatabaseRequest) *serializer.Response {
	if req == nil || req.Conid == "" {
		return serializer.Fail(serializer.ParamsErr)
	}

	opened := dc.ensureOpened(req.Conid, req.Database)
	if opened != nil && opened.ServerVersion != nil {
		return serializer.SuccessData(serializer.SUCCESS, opened.ServerVersion)
	}
	return serializer.SuccessData(serializer.SUCCESS, nil)
}

func (dc *DatabaseConnections) Status(req *DatabaseRequest) *serializer.Response {
	existing := findByDatabaseConnection(dc.Opened, req.Conid, req.Database)

	if existing != nil {
		return serializer.SuccessData(serializer.SUCCESS, map[string]interface{}{
			"name":         existing.Status.Name,
			"message":      existing.Status.Message,
			"counter":      existing.Status.Counter,
			"analysedTime": existing.AnalysedTime,
		})
	}

	lastClosed := dc.Closed[fmt.Sprintf("%s/%s", req.Conid, req.Database)]
	if lastClosed != nil {
		return serializer.SuccessData(serializer.SUCCESS, map[string]interface{}{
			"analysedTime": lastClosed.AnalysedTime,
		})
	}
	return serializer.SuccessData(serializer.SUCCESS, map[string]string{
		"name":    "error",
		"message": "Not connected",
	})
}

func (dc *DatabaseConnections) sendRequest() {

}

type DatabaseKillRequest struct {
	databaseConnections
	Kill bool `json:"kill"`
}

func (dc *DatabaseConnections) close(conid, database string, kill bool) {
	existing := findByDatabaseConnection(dc.Opened, conid, database)
	if existing != nil {
		existing.Disconnected = true
		if kill {

		}
		dc.Opened = lo.Filter[*containers.OpenedDatabaseConnection](dc.Opened, func(item *containers.OpenedDatabaseConnection, _ int) bool {
			return item.Conid != conid || item.Database != database
		})

		dc.Closed[fmt.Sprintf("%s/%s", conid, database)] = &containers.DatabaseConnectionClosed{
			Structure:    existing.Structure,
			AnalysedTime: existing.AnalysedTime,
			Status: &containers.OpenedStatus{
				Name:    "error",
				Message: existing.Status.Message,
				Counter: existing.Status.Counter,
			},
		}

		utility.EmitChanged(Application.ctx, fmt.Sprintf("database-status-changed-%s-%s", conid, database))
	}
}

func (dc *DatabaseConnections) closeAll(conid string, kill bool) {
	list := lo.Filter[*containers.OpenedDatabaseConnection](dc.Opened, func(item *containers.OpenedDatabaseConnection, _ int) bool {
		return item.Conid == conid
	})

	for _, v := range list {
		dc.close(conid, v.Database, kill)
	}
}

func (dc *DatabaseConnections) Disconnect(req *DatabaseRequest) *serializer.Response {
	dc.close(req.Conid, req.Database, true)
	return serializer.SuccessData(serializer.SUCCESS, &containers.OpenedStatus{Name: "ok"})
}

func findByDatabaseConnection(s []*containers.OpenedDatabaseConnection, conid, database string) *containers.OpenedDatabaseConnection {
	existing, ok := lo.Find[*containers.OpenedDatabaseConnection](s, func(item *containers.OpenedDatabaseConnection) bool {
		return item != nil && item.Conid != "" && item.Conid == conid && item.Database != "" && item.Database == database
	})

	if existing != nil && ok {
		return existing
	}
	return nil
}

type SqlSelectRequest struct {
	databaseConnections
}

const sqlSelectResponse = `{"msgtype":"response","rows":[{"file_id":"1","file_former_name":"手机数码.jpg","file_name":"1641349689799523.jpg","file_size":"39210","user_id":"0","file_join_id":"0","file_join_type":0,"hash":"Fpz0dpAqNBOvDpWxLMYTbQh55Vu2","bucket":"shop-attach","status":3,"created_at":"2022-01-05 10:28:10","updated_at":"2022-01-05 10:28:13","expire_at":"2022-01-06 10:28:10"},{"file_id":"2","file_former_name":"手机通信.jpg","file_name":"1641349822230131.jpg","file_size":"18639","user_id":"0","file_join_id":"0","file_join_type":0,"hash":"FhXNYgX1Qoq3Zc-mchM_-kS7uBQJ","bucket":"shop-attach","status":3,"created_at":"2022-01-05 10:30:23","updated_at":"2022-01-05 10:30:28","expire_at":"2022-01-06 10:30:23"},{"file_id":"3","file_former_name":"智能设备.jpg","file_name":"1641349866733432.jpg","file_size":"15691","user_id":"0","file_join_id":"0","file_join_type":0,"hash":"FguNw9ix6gpDhzRY8LXkefIKq3m4","bucket":"shop-attach","status":3,"created_at":"2022-01-05 10:31:06","updated_at":"2022-01-05 10:31:09","expire_at":"2022-01-06 10:31:06"},{"file_id":"4","file_former_name":"珠宝钟表.jpg","file_name":"1641349918485103.jpg","file_size":"39843","user_id":"0","file_join_id":"0","file_join_type":0,"hash":"Fs0s0lHAwmCq8Po7uStauM4d7iUK","bucket":"shop-attach","status":3,"created_at":"2022-01-05 10:31:59","updated_at":"2022-01-05 10:32:01","expire_at":"2022-01-06 10:31:59"},{"file_id":"5","file_former_name":"美妆护肤.jpg","file_name":"1641350062303094.jpg","file_size":"55132","user_id":"0","file_join_id":"0","file_join_type":0,"hash":"FkBkE6g4ovLpD-TYHSq1r8_UJp0-","bucket":"shop-attach","status":3,"created_at":"2022-01-05 10:34:22","updated_at":"2022-01-05 10:34:24","expire_at":"2022-01-06 10:34:22"},{"file_id":"6","file_former_name":"运动服饰.jpg","file_name":"1641350090842956.jpg","file_size":"55643","user_id":"0","file_join_id":"0","file_join_type":0,"hash":"FiPeHFO9ABeFRJoxa9YyUZpSLq9q","bucket":"shop-attach","status":3,"created_at":"2022-01-05 10:34:50","updated_at":"2022-01-05 10:34:54","expire_at":"2022-01-06 10:34:50"},{"file_id":"7","file_former_name":"美味零食.jpg","file_name":"1641350177654581.jpg","file_size":"43265","user_id":"0","file_join_id":"0","file_join_type":0,"hash":"FjyxG-nj_cIOtZCPSZZOVA3m4Efn","bucket":"shop-attach","status":3,"created_at":"2022-01-05 10:36:18","updated_at":"2022-01-05 10:36:19","expire_at":"2022-01-06 10:36:18"}],"columns":[{"columnName":"file_id"},{"columnName":"file_former_name"},{"columnName":"file_name"},{"columnName":"file_size"},{"columnName":"user_id"},{"columnName":"file_join_id"},{"columnName":"file_join_type"},{"columnName":"hash"},{"columnName":"bucket"},{"columnName":"status"},{"columnName":"created_at"},{"columnName":"updated_at"},{"columnName":"expire_at"}]}`

func (dc *DatabaseConnections) SqlSelect(req *SqlSelectRequest) *serializer.Response {
	unmarshal, err := utility.JsonUnmarshal([]byte(sqlSelectResponse))
	if err != nil {
		return serializer.Fail(err.Error())
	} else {
		return serializer.SuccessData(serializer.SUCCESS, unmarshal)
	}
}

func sendRequest(conn string) {
	uuid.NewV4().String()
}
