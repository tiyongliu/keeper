package bridge

import (
	"fmt"
	"github.com/samber/lo"
	"keeper/app/db/standard/modules"
	"keeper/app/pkg/containers"
	"keeper/app/pkg/logger"
	"keeper/app/pkg/serializer"
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

func (dc *DatabaseConnections) handleVersion(conid, database string, version *modules.Version) {
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
	dc.DatabaseConnection.ResetVars()

	go dc.DatabaseConnection.Connect(ch, newOpened, structure)
	go dc.consumer(ch, conid, database)
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

func (dc *DatabaseConnections) consumer(chData <-chan *containers.EchoMessage, conid, database string) {
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
				dc.handleVersion(conid, database, message.Payload.(*modules.Version))
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

func (dc *DatabaseConnections) sendRequest(conn *containers.OpenedDatabaseConnection, message *containers.EchoMessage) *containers.EchoMessage {
	if message == nil {
		return nil
	}
	switch message.MsgType {
	case "sqlSelect":
		return dc.DatabaseConnection.HandleSqlSelect(Application.ctx, conn, message.Payload)
	default:
		return nil
	}
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
	Select interface{}
}

const sqlSelectResponse1 = `{"msgtype":"response","rows":[{"att_id":"1","user_id":"1","real_name":"4b4a48ace185c67ac72ccb185821c5a9.jpg","name":"388393d9-7f10-4e4e-bf3d-698c2177fd38-wg2yz.jpeg","att_dir":"2022/02/08/","att_size":"55008","att_type":"jpeg","content_type":"image/jpeg","pid":"4","image_type":"1","key":"","hash":"","bucket":"","status":3,"created_at":"2022-02-08 11:52:38","updated_at":"2022-02-08 11:52:38"},{"att_id":"2","user_id":"1","real_name":"7edecfab5b2238b71e1520846fdf0918.jpg","name":"372cc0c6-b2d0-4e97-a730-2689c36391c5-ytcek.jpeg","att_dir":"2022/02/08/","att_size":"12442","att_type":"jpeg","content_type":"image/jpeg","pid":"4","image_type":"1","key":"","hash":"","bucket":"","status":3,"created_at":"2022-02-08 11:52:38","updated_at":"2022-02-08 11:52:38"},{"att_id":"3","user_id":"1","real_name":"05d25b7ed0675129acb3ed5ed9ea2baf.jpg","name":"10dbc307-23ee-40a1-99af-f7190e86e677-zhd74.jpeg","att_dir":"2022/02/08/","att_size":"165814","att_type":"jpeg","content_type":"image/jpeg","pid":"4","image_type":"1","key":"","hash":"","bucket":"","status":3,"created_at":"2022-02-08 11:52:38","updated_at":"2022-02-08 11:52:38"},{"att_id":"4","user_id":"1","real_name":"1ad625eb9de955df4114d64c20cb80e9.jpg","name":"337eaf7d-5394-4e17-898e-aeb77afbef3a-wg2yz.jpeg","att_dir":"2022/02/08/","att_size":"127406","att_type":"jpeg","content_type":"image/jpeg","pid":"4","image_type":"1","key":"","hash":"","bucket":"","status":3,"created_at":"2022-02-08 11:52:38","updated_at":"2022-02-08 11:52:38"},{"att_id":"5","user_id":"1","real_name":"1c2a3ecf19b53b749dd26269917a3fc6.jpg","name":"0726024b-2116-4071-8de3-d7240f662aa4-zhd74.jpeg","att_dir":"2022/02/08/","att_size":"299543","att_type":"jpeg","content_type":"image/jpeg","pid":"4","image_type":"1","key":"","hash":"","bucket":"","status":3,"created_at":"2022-02-08 11:52:38","updated_at":"2022-02-08 11:52:38"},{"att_id":"6","user_id":"1","real_name":"6fecaeb973533cec6c13b0581429b73f.jpg","name":"8b3dbccb-18bb-4ce1-b3d7-3243583e7fd1-g89qz.jpeg","att_dir":"2022/02/08/","att_size":"599667","att_type":"jpeg","content_type":"image/jpeg","pid":"4","image_type":"1","key":"","hash":"","bucket":"","status":3,"created_at":"2022-02-08 11:52:38","updated_at":"2022-02-08 11:52:38"},{"att_id":"7","user_id":"1","real_name":"9fb58503231b4a84a781a08d76c5f858.jpg","name":"d37ea8b2-ee3d-4c1e-915f-e4cdee1fda19-2ti5n.jpeg","att_dir":"2022/02/08/","att_size":"353459","att_type":"jpeg","content_type":"image/jpeg","pid":"4","image_type":"1","key":"","hash":"","bucket":"","status":3,"created_at":"2022-02-08 11:52:39","updated_at":"2022-02-08 11:52:39"},{"att_id":"8","user_id":"1","real_name":"59f89022fbd54af480f157cc4b69f689.jpg","name":"255d53d3-fb5d-4544-b8f9-d087483e48fd-pcvp2.jpeg","att_dir":"2022/02/08/","att_size":"563381","att_type":"jpeg","content_type":"image/jpeg","pid":"4","image_type":"1","key":"","hash":"","bucket":"","status":3,"created_at":"2022-02-08 11:52:39","updated_at":"2022-02-08 11:52:39"},{"att_id":"9","user_id":"1","real_name":"429e9fa99080f76ce67e7a3d91ed3f1c.jpg","name":"b2f7268a-b5f5-4895-afa5-cee3a59a026c-j6nvg.jpeg","att_dir":"2022/02/08/","att_size":"21471","att_type":"jpeg","content_type":"image/jpeg","pid":"4","image_type":"1","key":"","hash":"","bucket":"","status":3,"created_at":"2022-02-08 11:52:39","updated_at":"2022-02-08 11:52:39"},{"att_id":"10","user_id":"1","real_name":"73c67be46068dcc3f0acd2157014742b.jpg","name":"e5df7af7-ef26-44a8-a61a-ee17d96c0ea3-msmov.jpeg","att_dir":"2022/02/08/","att_size":"183220","att_type":"jpeg","content_type":"image/jpeg","pid":"4","image_type":"1","key":"","hash":"","bucket":"","status":3,"created_at":"2022-02-08 11:52:39","updated_at":"2022-02-08 11:52:39"},{"att_id":"11","user_id":"1","real_name":"9152b56db573d91d32649f14704c628e.jpg","name":"68378cb7-c113-4c5b-bb56-4d7f8e17a233-vc470.jpeg","att_dir":"2022/02/08/","att_size":"454962","att_type":"jpeg","content_type":"image/jpeg","pid":"4","image_type":"1","key":"","hash":"","bucket":"","status":3,"created_at":"2022-02-08 11:52:39","updated_at":"2022-02-08 11:52:39"},{"att_id":"12","user_id":"1","real_name":"98373ea5a76782a789d75871f8bf9d35.jpg","name":"cc288b91-614e-4ee5-96b1-8145ed49f943-3di3b.jpeg","att_dir":"2022/02/08/","att_size":"47619","att_type":"jpeg","content_type":"image/jpeg","pid":"4","image_type":"1","key":"","hash":"","bucket":"","status":3,"created_at":"2022-02-08 11:52:39","updated_at":"2022-02-08 11:52:39"},{"att_id":"13","user_id":"1","real_name":"97800f14d4241e6c6c29509e7c729780.jpg","name":"230ad383-bbee-425e-b383-b1fd827dea48-3di3b.jpeg","att_dir":"2022/02/08/","att_size":"222420","att_type":"jpeg","content_type":"image/jpeg","pid":"4","image_type":"1","key":"","hash":"","bucket":"","status":3,"created_at":"2022-02-08 11:52:39","updated_at":"2022-02-08 11:52:39"},{"att_id":"14","user_id":"1","real_name":"31485c1e2abc8b2b3bce4a3f151a6ff1.jpg","name":"63eb2d82-2e5c-417c-93a2-1d823d90952d-o4sml.jpeg","att_dir":"2022/02/08/","att_size":"553555","att_type":"jpeg","content_type":"image/jpeg","pid":"4","image_type":"1","key":"","hash":"","bucket":"","status":3,"created_at":"2022-02-08 11:52:39","updated_at":"2022-02-08 11:52:39"},{"att_id":"15","user_id":"1","real_name":"0155529b543106bd204687142f7eea67.jpg","name":"fb0d500d-ea58-4ccc-befa-828ea1fcd7e0-o4sml.jpeg","att_dir":"2022/02/08/","att_size":"421603","att_type":"jpeg","content_type":"image/jpeg","pid":"4","image_type":"1","key":"","hash":"","bucket":"","status":3,"created_at":"2022-02-08 11:52:39","updated_at":"2022-02-08 11:52:39"},{"att_id":"16","user_id":"1","real_name":"a9247bcc28ab56a89b47573d2beb4f44.jpg","name":"dcba8dad-ff63-4ecf-9701-0f13e541c64b-6rt67.jpeg","att_dir":"2022/02/08/","att_size":"20213","att_type":"jpeg","content_type":"image/jpeg","pid":"4","image_type":"1","key":"","hash":"","bucket":"","status":3,"created_at":"2022-02-08 11:52:39","updated_at":"2022-02-08 11:52:39"},{"att_id":"17","user_id":"1","real_name":"d6a37a6c6bbb5fdaf5a3aa98788df789.jpg","name":"0446887c-86f9-4abe-90e8-34c56ca4bf85-24tsq.jpeg","att_dir":"2022/02/08/","att_size":"334372","att_type":"jpeg","content_type":"image/jpeg","pid":"4","image_type":"1","key":"","hash":"","bucket":"","status":3,"created_at":"2022-02-08 11:52:39","updated_at":"2022-02-08 11:52:39"},{"att_id":"18","user_id":"1","real_name":"d60b5551fb88586dd6d009a4c3db7823.jpg","name":"a83c75fe-9745-48f8-a05f-d2e2994d0390-pcdx0.jpeg","att_dir":"2022/02/08/","att_size":"17073","att_type":"jpeg","content_type":"image/jpeg","pid":"4","image_type":"1","key":"","hash":"","bucket":"","status":3,"created_at":"2022-02-08 11:52:39","updated_at":"2022-02-08 11:52:39"},{"att_id":"19","user_id":"1","real_name":"d72acdaa04def74b677a9fe98c8f2364.jpg","name":"66b13e64-a006-42b1-91c6-ce7f21632d49-zft1h.jpeg","att_dir":"2022/02/08/","att_size":"83334","att_type":"jpeg","content_type":"image/jpeg","pid":"4","image_type":"1","key":"","hash":"","bucket":"","status":3,"created_at":"2022-02-08 11:52:39","updated_at":"2022-02-08 11:52:39"},{"att_id":"20","user_id":"1","real_name":"e0d8fa45671b5927872f9cb68c8c0c56.jpg","name":"473e82b2-8968-4e7e-8599-5d7d2966cec4-t8fvk.jpeg","att_dir":"2022/02/08/","att_size":"178404","att_type":"jpeg","content_type":"image/jpeg","pid":"4","image_type":"1","key":"","hash":"","bucket":"","status":3,"created_at":"2022-02-08 11:52:39","updated_at":"2022-02-08 11:52:39"},{"att_id":"21","user_id":"1","real_name":"bef0711f61bb1f8c9599ef7509e4b7ba.jpg","name":"fdadad8c-6905-4e98-a7e5-aaf4b44f31a8-fylui.jpeg","att_dir":"2022/02/08/","att_size":"82861","att_type":"jpeg","content_type":"image/jpeg","pid":"4","image_type":"1","key":"","hash":"","bucket":"","status":3,"created_at":"2022-02-08 11:52:54","updated_at":"2022-02-08 11:52:54"},{"att_id":"22","user_id":"1","real_name":"b38aef91cb8be2eb75e5579158a13344.jpg","name":"e92bcf11-1c6b-49c7-a45a-f32f5ee06fff-fylui.jpeg","att_dir":"2022/02/08/","att_size":"221237","att_type":"jpeg","content_type":"image/jpeg","pid":"4","image_type":"1","key":"","hash":"","bucket":"","status":3,"created_at":"2022-02-08 11:52:54","updated_at":"2022-02-08 11:52:54"},{"att_id":"23","user_id":"1","real_name":"89750263349d7aabb1380c7d0a0ce11b.jpg","name":"705d74c7-9aed-48d1-a76e-43e9187f1bf9-vgrfl.jpeg","att_dir":"2022/02/08/","att_size":"613173","att_type":"jpeg","content_type":"image/jpeg","pid":"4","image_type":"1","key":"","hash":"","bucket":"","status":3,"created_at":"2022-02-08 11:52:54","updated_at":"2022-02-08 11:52:54"},{"att_id":"24","user_id":"1","real_name":"c246e78dc1a6240c79abfc3d5ef411ca.jpg","name":"334b8e06-e83c-482b-a7ab-5347e27a9486-qhzvc.jpeg","att_dir":"2022/02/08/","att_size":"199681","att_type":"jpeg","content_type":"image/jpeg","pid":"4","image_type":"1","key":"","hash":"","bucket":"","status":3,"created_at":"2022-02-08 11:52:54","updated_at":"2022-02-08 11:52:54"},{"att_id":"25","user_id":"1","real_name":"bf9a9857a14c1f67816b2ed11e8f413e.jpg","name":"eedfe8b8-65ae-4157-b4f5-09ea7751f0b4-qhzvc.jpeg","att_dir":"2022/02/08/","att_size":"255865","att_type":"jpeg","content_type":"image/jpeg","pid":"4","image_type":"1","key":"","hash":"","bucket":"","status":3,"created_at":"2022-02-08 11:52:54","updated_at":"2022-02-08 11:52:54"},{"att_id":"26","user_id":"1","real_name":"c8322dcecd2843856789d5cea9646340.jpg","name":"2e6b0fe8-b057-4020-a75c-4ca90a523f43-yzdea.jpeg","att_dir":"2022/02/08/","att_size":"108274","att_type":"jpeg","content_type":"image/jpeg","pid":"4","image_type":"1","key":"","hash":"","bucket":"","status":3,"created_at":"2022-02-08 11:52:54","updated_at":"2022-02-08 11:52:54"},{"att_id":"27","user_id":"1","real_name":"e0ccd5a9103909277f328747b40c5ccf.jpg","name":"6ff77fbe-db70-4312-8189-c262c15a3b7f-sskfs.jpeg","att_dir":"2022/02/08/","att_size":"380281","att_type":"jpeg","content_type":"image/jpeg","pid":"4","image_type":"1","key":"","hash":"","bucket":"","status":3,"created_at":"2022-02-08 11:52:54","updated_at":"2022-02-08 11:52:54"},{"att_id":"28","user_id":"1","real_name":"e2fd950c316c7ded99baa5772e5488eb.png","name":"473a0890-513d-4172-97fc-f3b8434b4431-o1vrd.png","att_dir":"2022/02/08/","att_size":"158162","att_type":"png","content_type":"image/png","pid":"4","image_type":"1","key":"","hash":"","bucket":"","status":3,"created_at":"2022-02-08 11:52:54","updated_at":"2022-02-08 11:52:54"},{"att_id":"29","user_id":"1","real_name":"e5dadd32dbfd0f7762763c6d048f4dd3.jpg","name":"8e13203e-c44a-49b4-9cd8-8638caacd075-2nmqq.jpeg","att_dir":"2022/02/08/","att_size":"189776","att_type":"jpeg","content_type":"image/jpeg","pid":"4","image_type":"1","key":"","hash":"","bucket":"","status":3,"created_at":"2022-02-08 11:52:54","updated_at":"2022-02-08 11:52:54"},{"att_id":"30","user_id":"1","real_name":"e53a6d40efb95141291de1211dce1682.jpg","name":"2edc6781-bc15-40b7-86d3-a8cdf9e3f244-x83ns.jpeg","att_dir":"2022/02/08/","att_size":"128194","att_type":"jpeg","content_type":"image/jpeg","pid":"4","image_type":"1","key":"","hash":"","bucket":"","status":3,"created_at":"2022-02-08 11:52:54","updated_at":"2022-02-08 11:52:54"},{"att_id":"31","user_id":"1","real_name":"f73a2585658c3889e9f92b68386d8a26.jpg","name":"755a9bbf-97ab-4aeb-aa47-ab1d86504c68-m8f5y.jpeg","att_dir":"2022/02/08/","att_size":"207775","att_type":"jpeg","content_type":"image/jpeg","pid":"4","image_type":"1","key":"","hash":"","bucket":"","status":3,"created_at":"2022-02-08 11:52:54","updated_at":"2022-02-08 11:52:54"},{"att_id":"32","user_id":"1","real_name":"f8873f426195ab9e287bdaad0de46638.jpg","name":"352db06a-c3d2-4145-be4e-1022b1d83ed2-5ojj7.jpeg","att_dir":"2022/02/08/","att_size":"450604","att_type":"jpeg","content_type":"image/jpeg","pid":"4","image_type":"1","key":"","hash":"","bucket":"","status":3,"created_at":"2022-02-08 11:52:54","updated_at":"2022-02-08 11:52:54"},{"att_id":"33","user_id":"1","real_name":"0e5043730ab5ea9fbec861a9cb496c27.png","name":"800601aa-4147-4586-bce8-9a1ad767b1f9-l2sr1.png","att_dir":"2022/02/08/","att_size":"11005","att_type":"png","content_type":"image/png","pid":"2","image_type":"1","key":"","hash":"","bucket":"","status":3,"created_at":"2022-02-08 11:54:17","updated_at":"2022-02-08 11:54:17"},{"att_id":"34","user_id":"1","real_name":"2bb068405ddda07a4c719bf9bd7812d8.jpg","name":"f1007471-f009-4a83-bce9-8e8d6d1ba70f-8r6vb.jpeg","att_dir":"2022/02/08/","att_size":"3604","att_type":"jpeg","content_type":"image/jpeg","pid":"2","image_type":"1","key":"","hash":"","bucket":"","status":3,"created_at":"2022-02-08 11:54:17","updated_at":"2022-02-08 11:54:17"},{"att_id":"35","user_id":"1","real_name":"7ef129224510e00d9b0f9da1057c6ab2.png","name":"34f1b914-1d36-4933-9436-b30f6704f9fd-mdn3k.png","att_dir":"2022/02/08/","att_size":"13384","att_type":"png","content_type":"image/png","pid":"2","image_type":"1","key":"","hash":"","bucket":"","status":3,"created_at":"2022-02-08 11:54:17","updated_at":"2022-02-08 11:54:17"},{"att_id":"36","user_id":"1","real_name":"8d58a46bd06e65a0a09dac36b919edae.png","name":"47781da7-e68a-4420-ac67-d01993541a0d-462r0.png","att_dir":"2022/02/08/","att_size":"8244","att_type":"png","content_type":"image/png","pid":"2","image_type":"1","key":"","hash":"","bucket":"","status":3,"created_at":"2022-02-08 11:54:17","updated_at":"2022-02-08 11:54:17"},{"att_id":"37","user_id":"1","real_name":"915e3db85e4ccdaf56409141ec3233db.png","name":"09589af9-967f-4f94-ba66-d146da7f8b1e-2y6yd.png","att_dir":"2022/02/08/","att_size":"11612","att_type":"png","content_type":"image/png","pid":"2","image_type":"1","key":"","hash":"","bucket":"","status":3,"created_at":"2022-02-08 11:54:17","updated_at":"2022-02-08 11:54:17"},{"att_id":"38","user_id":"1","real_name":"3459e11bd015ccac1773d61a8e086164.png","name":"f7cb2731-9855-4ad9-825e-dd9a9970fd0b-2y6yd.png","att_dir":"2022/02/08/","att_size":"16174","att_type":"png","content_type":"image/png","pid":"2","image_type":"1","key":"","hash":"","bucket":"","status":3,"created_at":"2022-02-08 11:54:17","updated_at":"2022-02-08 11:54:17"},{"att_id":"39","user_id":"1","real_name":"a195c37016d14ec32e70c7d59bafd470.png","name":"efae8368-e05d-48b5-9edf-3010c7537727-hi9ws.png","att_dir":"2022/02/08/","att_size":"13374","att_type":"png","content_type":"image/png","pid":"2","image_type":"1","key":"","hash":"","bucket":"","status":3,"created_at":"2022-02-08 11:54:18","updated_at":"2022-02-08 11:54:18"},{"att_id":"40","user_id":"1","real_name":"bb5bfda56b7564fbe59a8c3001e1a957.png","name":"ffae262b-dd44-4913-9ad9-28d8e8b0961a-s7pzs.png","att_dir":"2022/02/08/","att_size":"10304","att_type":"png","content_type":"image/png","pid":"2","image_type":"1","key":"","hash":"","bucket":"","status":3,"created_at":"2022-02-08 11:54:18","updated_at":"2022-02-08 11:54:18"},{"att_id":"41","user_id":"1","real_name":"bcb10769df9f3bdd2d0e189801934862.png","name":"678cc57d-2ddb-4b53-8b51-8a8b03229831-nekam.png","att_dir":"2022/02/08/","att_size":"14584","att_type":"png","content_type":"image/png","pid":"2","image_type":"1","key":"","hash":"","bucket":"","status":3,"created_at":"2022-02-08 11:54:18","updated_at":"2022-02-08 11:54:18"},{"att_id":"42","user_id":"1","real_name":"fbd6ad64bb79c42634628fbb336f8d46.png","name":"7743039f-1bc8-4ea6-9546-eb3cfa441864-zx1t4.png","att_dir":"2022/02/08/","att_size":"17361","att_type":"png","content_type":"image/png","pid":"2","image_type":"1","key":"","hash":"","bucket":"","status":3,"created_at":"2022-02-08 11:54:18","updated_at":"2022-02-08 11:54:18"},{"att_id":"43","user_id":"1","real_name":"cd557f6a6068fdc17c85d2c04cb4a3ee.jpg","name":"adb909d4-868d-4017-898d-5d2c4cbf3fb7-zx1t4.jpeg","att_dir":"2022/02/08/","att_size":"26852","att_type":"jpeg","content_type":"image/jpeg","pid":"2","image_type":"1","key":"","hash":"","bucket":"","status":3,"created_at":"2022-02-08 11:54:18","updated_at":"2022-02-08 11:54:18"},{"att_id":"44","user_id":"1","real_name":"f5115ff4d97804377d6ff44d1dc6fa1c.jpg","name":"f50056c1-5680-41a3-8a55-f930332a9e5e-otpuv.jpeg","att_dir":"2022/02/08/","att_size":"435312","att_type":"jpeg","content_type":"image/jpeg","pid":"2","image_type":"1","key":"","hash":"","bucket":"","status":3,"created_at":"2022-02-08 11:54:18","updated_at":"2022-02-08 11:54:18"},{"att_id":"45","user_id":"1","real_name":"d9b193d7c3ed2793508977fe865fc945.jpeg","name":"248906d8-332f-4791-b8c3-74863e22a476-l3dn4.jpeg","att_dir":"2022/02/08/","att_size":"165730","att_type":"jpeg","content_type":"image/jpeg","pid":"3","image_type":"1","key":"","hash":"","bucket":"","status":3,"created_at":"2022-02-08 11:54:32","updated_at":"2022-02-08 11:54:32"},{"att_id":"46","user_id":"1","real_name":"e973a232df7dc7c3f421c1116ce958e2.jpeg","name":"f13be698-2194-4683-8526-ee35d33d89dd-iw51r.jpeg","att_dir":"2022/02/08/","att_size":"97039","att_type":"jpeg","content_type":"image/jpeg","pid":"3","image_type":"1","key":"","hash":"","bucket":"","status":3,"created_at":"2022-02-08 11:54:32","updated_at":"2022-02-08 11:54:32"},{"att_id":"47","user_id":"1","real_name":"0ce3c013963d2554edd228bbe7299058.jpg","name":"5d69e987-4f74-4408-a0f8-93e1bdc84902-ip2qe.jpeg","att_dir":"2022/02/08/","att_size":"111069","att_type":"jpeg","content_type":"image/jpeg","pid":"5","image_type":"1","key":"","hash":"","bucket":"","status":3,"created_at":"2022-02-08 11:55:10","updated_at":"2022-02-08 11:55:10"},{"att_id":"48","user_id":"1","real_name":"0ababb301df78c1858a3f77887624cba.jpg","name":"1c34e186-5ea3-4829-b314-0eb40ca412e1-dlpwl.jpeg","att_dir":"2022/02/08/","att_size":"318986","att_type":"jpeg","content_type":"image/jpeg","pid":"5","image_type":"1","key":"","hash":"","bucket":"","status":3,"created_at":"2022-02-08 11:55:10","updated_at":"2022-02-08 11:55:10"},{"att_id":"49","user_id":"1","real_name":"1c42b2a4d7fe05e108eda936fd494235 (1).jpg","name":"8438cdf8-78d0-44d9-97e3-916c7b5472ce-ffvm0.jpeg","att_dir":"2022/02/08/","att_size":"182430","att_type":"jpeg","content_type":"image/jpeg","pid":"5","image_type":"1","key":"","hash":"","bucket":"","status":3,"created_at":"2022-02-08 11:55:10","updated_at":"2022-02-08 11:55:10"},{"att_id":"50","user_id":"1","real_name":"1c42b2a4d7fe05e108eda936fd494235.jpg","name":"52406289-e411-44b7-a992-455e889a92db-3alht.jpeg","att_dir":"2022/02/08/","att_size":"182430","att_type":"jpeg","content_type":"image/jpeg","pid":"5","image_type":"1","key":"","hash":"","bucket":"","status":3,"created_at":"2022-02-08 11:55:10","updated_at":"2022-02-08 11:55:10"},{"att_id":"51","user_id":"1","real_name":"2e8b4da1e5b466c0ffd4ebce11eba549 (1).jpg","name":"acbc2e08-c5d4-4cfc-9c50-a70bd59cc632-j4m5r.jpeg","att_dir":"2022/02/08/","att_size":"525861","att_type":"jpeg","content_type":"image/jpeg","pid":"5","image_type":"1","key":"","hash":"","bucket":"","status":3,"created_at":"2022-02-08 11:55:10","updated_at":"2022-02-08 11:55:10"},{"att_id":"52","user_id":"1","real_name":"2e8b4da1e5b466c0ffd4ebce11eba549 (2).jpg","name":"7971ccd8-005c-4908-a581-fd7068e3fcd5-j4m5r.jpeg","att_dir":"2022/02/08/","att_size":"525861","att_type":"jpeg","content_type":"image/jpeg","pid":"5","image_type":"1","key":"","hash":"","bucket":"","status":3,"created_at":"2022-02-08 11:55:10","updated_at":"2022-02-08 11:55:10"},{"att_id":"53","user_id":"1","real_name":"2e8b4da1e5b466c0ffd4ebce11eba549.jpg","name":"e13beede-6285-422d-ba09-731847391676-974pj.jpeg","att_dir":"2022/02/08/","att_size":"525861","att_type":"jpeg","content_type":"image/jpeg","pid":"5","image_type":"1","key":"","hash":"","bucket":"","status":3,"created_at":"2022-02-08 11:55:10","updated_at":"2022-02-08 11:55:10"},{"att_id":"54","user_id":"1","real_name":"6bc02547d04a348666164437fa754d9d.jpg","name":"b65d90e9-4238-4b8d-a192-4ab8449783c1-dkqq8.jpeg","att_dir":"2022/02/08/","att_size":"308098","att_type":"jpeg","content_type":"image/jpeg","pid":"5","image_type":"1","key":"","hash":"","bucket":"","status":3,"created_at":"2022-02-08 11:55:10","updated_at":"2022-02-08 11:55:10"},{"att_id":"55","user_id":"1","real_name":"6d4cfe2e8c4990df6549f4f202c0fe92.jpg","name":"46775e47-053c-4f5e-964f-b88c06196180-yiihh.jpeg","att_dir":"2022/02/08/","att_size":"338723","att_type":"jpeg","content_type":"image/jpeg","pid":"5","image_type":"1","key":"","hash":"","bucket":"","status":3,"created_at":"2022-02-08 11:55:10","updated_at":"2022-02-08 11:55:10"},{"att_id":"56","user_id":"1","real_name":"97a39bc9c640c93a86abdd2d6777850c.jpg","name":"74c1a5b2-bfe1-4efb-b09f-96cfe27053fd-99q8h.jpeg","att_dir":"2022/02/08/","att_size":"76716","att_type":"jpeg","content_type":"image/jpeg","pid":"5","image_type":"1","key":"","hash":"","bucket":"","status":3,"created_at":"2022-02-08 11:55:10","updated_at":"2022-02-08 11:55:10"},{"att_id":"57","user_id":"1","real_name":"410d7320694eef14b999e3a4e9807c93 (2).jpg","name":"b784f5d2-4492-4a34-bd9a-9c259a7fc4fb-lx4jv.jpeg","att_dir":"2022/02/08/","att_size":"648620","att_type":"jpeg","content_type":"image/jpeg","pid":"5","image_type":"1","key":"","hash":"","bucket":"","status":3,"created_at":"2022-02-08 11:55:10","updated_at":"2022-02-08 11:55:10"},{"att_id":"58","user_id":"1","real_name":"410d7320694eef14b999e3a4e9807c93 (1).jpg","name":"43b40d0b-a726-4b99-9575-3ea05534f78b-lx4jv.jpeg","att_dir":"2022/02/08/","att_size":"648620","att_type":"jpeg","content_type":"image/jpeg","pid":"5","image_type":"1","key":"","hash":"","bucket":"","status":3,"created_at":"2022-02-08 11:55:10","updated_at":"2022-02-08 11:55:10"},{"att_id":"59","user_id":"1","real_name":"410d7320694eef14b999e3a4e9807c93.jpg","name":"299b6d8f-03f7-4d37-a6df-849520fd8c60-kb8vq.jpeg","att_dir":"2022/02/08/","att_size":"648620","att_type":"jpeg","content_type":"image/jpeg","pid":"5","image_type":"1","key":"","hash":"","bucket":"","status":3,"created_at":"2022-02-08 11:55:10","updated_at":"2022-02-08 11:55:10"},{"att_id":"60","user_id":"1","real_name":"541c9a533db225261560ecd1e5577590 (1).jpg","name":"0dc7681d-7bd0-485e-943c-fa1b4cdb7a78-44o1t.jpeg","att_dir":"2022/02/08/","att_size":"574306","att_type":"jpeg","content_type":"image/jpeg","pid":"5","image_type":"1","key":"","hash":"","bucket":"","status":3,"created_at":"2022-02-08 11:55:10","updated_at":"2022-02-08 11:55:10"},{"att_id":"61","user_id":"1","real_name":"541c9a533db225261560ecd1e5577590 (2).jpg","name":"4b76c7bb-7e99-4391-8472-4d91c9f9d0a7-k7lcz.jpeg","att_dir":"2022/02/08/","att_size":"574306","att_type":"jpeg","content_type":"image/jpeg","pid":"5","image_type":"1","key":"","hash":"","bucket":"","status":3,"created_at":"2022-02-08 11:55:10","updated_at":"2022-02-08 11:55:10"},{"att_id":"62","user_id":"1","real_name":"701f0fb8468c5c7f626aef181e477d1c (1).jpg","name":"d6514d03-f985-4192-8429-eb559cb1c2a2-fft8n.jpeg","att_dir":"2022/02/08/","att_size":"172855","att_type":"jpeg","content_type":"image/jpeg","pid":"5","image_type":"1","key":"","hash":"","bucket":"","status":3,"created_at":"2022-02-08 11:55:10","updated_at":"2022-02-08 11:55:10"},{"att_id":"63","user_id":"1","real_name":"883c0ee6ba484b8982005a1402d14e33 (1).jpg","name":"16a8f5fa-fb80-4d04-aae3-5d9818b3ae0a-0gv42.jpeg","att_dir":"2022/02/08/","att_size":"427723","att_type":"jpeg","content_type":"image/jpeg","pid":"5","image_type":"1","key":"","hash":"","bucket":"","status":3,"created_at":"2022-02-08 11:55:10","updated_at":"2022-02-08 11:55:10"},{"att_id":"64","user_id":"1","real_name":"b3c47f50bf6b2efef4dbe62f61eba800 (5).jpg","name":"d04e0575-5a78-454d-99b6-ff6934057289-p03n5.jpeg","att_dir":"2022/02/08/","att_size":"672135","att_type":"jpeg","content_type":"image/jpeg","pid":"5","image_type":"1","key":"","hash":"","bucket":"","status":3,"created_at":"2022-02-08 11:55:11","updated_at":"2022-02-08 11:55:11"},{"att_id":"65","user_id":"1","real_name":"c2b36b28fa1a15a0a3cd975ce6e69cc5.jpg","name":"8b1601a6-ff62-4ae3-9003-9630010af207-6x6d4.jpeg","att_dir":"2022/02/08/","att_size":"271327","att_type":"jpeg","content_type":"image/jpeg","pid":"5","image_type":"1","key":"","hash":"","bucket":"","status":3,"created_at":"2022-02-08 11:55:11","updated_at":"2022-02-08 11:55:11"},{"att_id":"66","user_id":"1","real_name":"b3c47f50bf6b2efef4dbe62f61eba800.jpg","name":"637c90ab-f3a3-4781-ac92-2e9207f8ff77-4p0ts.jpeg","att_dir":"2022/02/08/","att_size":"672135","att_type":"jpeg","content_type":"image/jpeg","pid":"5","image_type":"1","key":"","hash":"","bucket":"","status":3,"created_at":"2022-02-08 11:55:11","updated_at":"2022-02-08 11:55:11"},{"att_id":"67","user_id":"1","real_name":"701f0fb8468c5c7f626aef181e477d1c.jpg","name":"3cdc7720-8ef3-4491-bbf5-bf1c7c2d2fc2-ro0df.jpeg","att_dir":"2022/02/08/","att_size":"172855","att_type":"jpeg","content_type":"image/jpeg","pid":"5","image_type":"1","key":"","hash":"","bucket":"","status":3,"created_at":"2022-02-08 11:55:16","updated_at":"2022-02-08 11:55:16"},{"att_id":"68","user_id":"1","real_name":"883c0ee6ba484b8982005a1402d14e33 (2).jpg","name":"dfcb8607-a5f8-43de-80d0-2b1795c18cc4-vd4zw.jpeg","att_dir":"2022/02/08/","att_size":"427723","att_type":"jpeg","content_type":"image/jpeg","pid":"5","image_type":"1","key":"","hash":"","bucket":"","status":3,"created_at":"2022-02-08 11:55:16","updated_at":"2022-02-08 11:55:16"},{"att_id":"69","user_id":"1","real_name":"0782de8538217e277e302a46aee53f49.jpg","name":"79024853-bc66-4700-83bd-24b75356d65d-mqdwt.jpeg","att_dir":"2022/02/08/","att_size":"54833","att_type":"jpeg","content_type":"image/jpeg","pid":"5","image_type":"1","key":"","hash":"","bucket":"","status":3,"created_at":"2022-02-08 11:55:16","updated_at":"2022-02-08 11:55:16"},{"att_id":"70","user_id":"1","real_name":"a0e9032b87864f869fef17b8d70e52d7.jpg","name":"3daef2e0-dc61-4e9a-b8c6-f3d237d88429-adclh.jpeg","att_dir":"2022/02/08/","att_size":"290498","att_type":"jpeg","content_type":"image/jpeg","pid":"5","image_type":"1","key":"","hash":"","bucket":"","status":3,"created_at":"2022-02-08 11:55:16","updated_at":"2022-02-08 11:55:16"},{"att_id":"71","user_id":"1","real_name":"883c0ee6ba484b8982005a1402d14e33.jpg","name":"983ad412-7a42-4955-82b6-be3e34d84588-ga1pg.jpeg","att_dir":"2022/02/08/","att_size":"427723","att_type":"jpeg","content_type":"image/jpeg","pid":"5","image_type":"1","key":"","hash":"","bucket":"","status":3,"created_at":"2022-02-08 11:55:16","updated_at":"2022-02-08 11:55:16"},{"att_id":"72","user_id":"1","real_name":"541c9a533db225261560ecd1e5577590.jpg","name":"78e68c8f-dce2-40ba-97dc-b028f5d61576-adclh.jpeg","att_dir":"2022/02/08/","att_size":"574306","att_type":"jpeg","content_type":"image/jpeg","pid":"5","image_type":"1","key":"","hash":"","bucket":"","status":3,"created_at":"2022-02-08 11:55:16","updated_at":"2022-02-08 11:55:16"},{"att_id":"73","user_id":"1","real_name":"b3c47f50bf6b2efef4dbe62f61eba800 (1).jpg","name":"4a2ea076-ea17-4835-8137-52ef5ee17869-lmefv.jpeg","att_dir":"2022/02/08/","att_size":"672135","att_type":"jpeg","content_type":"image/jpeg","pid":"5","image_type":"1","key":"","hash":"","bucket":"","status":3,"created_at":"2022-02-08 11:55:16","updated_at":"2022-02-08 11:55:16"},{"att_id":"74","user_id":"1","real_name":"b3c47f50bf6b2efef4dbe62f61eba800 (2).jpg","name":"e50bb71c-1833-40fb-ab42-6a603f92d5d8-jvmmg.jpeg","att_dir":"2022/02/08/","att_size":"672135","att_type":"jpeg","content_type":"image/jpeg","pid":"5","image_type":"1","key":"","hash":"","bucket":"","status":3,"created_at":"2022-02-08 11:55:16","updated_at":"2022-02-08 11:55:16"},{"att_id":"75","user_id":"1","real_name":"b9664e3c3f8c62d47e7b37df1336e968.jpg","name":"ed87a655-faef-44d1-a046-982f80a36e32-kleg8.jpeg","att_dir":"2022/02/08/","att_size":"51974","att_type":"jpeg","content_type":"image/jpeg","pid":"5","image_type":"1","key":"","hash":"","bucket":"","status":3,"created_at":"2022-02-08 11:55:16","updated_at":"2022-02-08 11:55:16"},{"att_id":"76","user_id":"1","real_name":"b3c47f50bf6b2efef4dbe62f61eba800 (4).jpg","name":"b3bd1a98-229d-4e7e-ab39-6665cc207d51-34rjf.jpeg","att_dir":"2022/02/08/","att_size":"672135","att_type":"jpeg","content_type":"image/jpeg","pid":"5","image_type":"1","key":"","hash":"","bucket":"","status":3,"created_at":"2022-02-08 11:55:16","updated_at":"2022-02-08 11:55:16"},{"att_id":"77","user_id":"1","real_name":"b3c47f50bf6b2efef4dbe62f61eba800 (3).jpg","name":"985a7595-6f7a-4433-976a-ddc75d378125-vq30z.jpeg","att_dir":"2022/02/08/","att_size":"672135","att_type":"jpeg","content_type":"image/jpeg","pid":"5","image_type":"1","key":"","hash":"","bucket":"","status":3,"created_at":"2022-02-08 11:55:16","updated_at":"2022-02-08 11:55:16"},{"att_id":"78","user_id":"1","real_name":"c0440cea35732278a0c718e25084e111.jpg","name":"c35cca87-4c49-4bc4-b48d-6bd0c1df75a1-0evda.jpeg","att_dir":"2022/02/08/","att_size":"317344","att_type":"jpeg","content_type":"image/jpeg","pid":"5","image_type":"1","key":"","hash":"","bucket":"","status":3,"created_at":"2022-02-08 11:55:16","updated_at":"2022-02-08 11:55:16"},{"att_id":"79","user_id":"1","real_name":"ce4c8e283c71e8ec7c114984dda29d04.jpg","name":"584e594e-062f-46b2-a37f-613a102914e7-u7ntb.jpeg","att_dir":"2022/02/08/","att_size":"324407","att_type":"jpeg","content_type":"image/jpeg","pid":"5","image_type":"1","key":"","hash":"","bucket":"","status":3,"created_at":"2022-02-08 11:55:16","updated_at":"2022-02-08 11:55:16"},{"att_id":"80","user_id":"1","real_name":"d9c83978fa94a2fc78a76ef2900186e0.jpg","name":"35298d44-6460-4cbd-8097-0e154acfd9e8-k0nkw.jpeg","att_dir":"2022/02/08/","att_size":"309151","att_type":"jpeg","content_type":"image/jpeg","pid":"5","image_type":"1","key":"","hash":"","bucket":"","status":3,"created_at":"2022-02-08 11:55:16","updated_at":"2022-02-08 11:55:16"},{"att_id":"81","user_id":"1","real_name":"361.jpg","name":"e7046272-b61d-411d-8d56-5895d67174e5-kwj53.jpeg","att_dir":"2022/02/08/","att_size":"8046","att_type":"jpeg","content_type":"image/jpeg","pid":"6","image_type":"1","key":"","hash":"","bucket":"","status":3,"created_at":"2022-02-08 13:54:50","updated_at":"2022-02-08 13:54:50"},{"att_id":"82","user_id":"1","real_name":"安踏.jpg","name":"0723b346-c02d-4e88-b736-240c44a030f6-b8lq7.jpeg","att_dir":"2022/02/08/","att_size":"10201","att_type":"jpeg","content_type":"image/jpeg","pid":"6","image_type":"1","key":"","hash":"","bucket":"","status":3,"created_at":"2022-02-08 13:54:50","updated_at":"2022-02-08 13:54:50"},{"att_id":"83","user_id":"1","real_name":"vans.png","name":"80983713-b556-4df5-abe5-a16361af6b1a-qri9s.png","att_dir":"2022/02/08/","att_size":"23281","att_type":"png","content_type":"image/png","pid":"6","image_type":"1","key":"","hash":"","bucket":"","status":3,"created_at":"2022-02-08 13:54:50","updated_at":"2022-02-08 13:54:50"},{"att_id":"84","user_id":"1","real_name":"阿迪达斯.jpg","name":"2acda7fb-5b08-4452-a759-24ef8daecdf8-qri9s.jpeg","att_dir":"2022/02/08/","att_size":"9234","att_type":"jpeg","content_type":"image/jpeg","pid":"6","image_type":"1","key":"","hash":"","bucket":"","status":3,"created_at":"2022-02-08 13:54:50","updated_at":"2022-02-08 13:54:50"},{"att_id":"85","user_id":"1","real_name":"斐乐.jpg","name":"73629b46-9a93-45a9-83dd-b418ec3a1a80-qri9s.jpeg","att_dir":"2022/02/08/","att_size":"5072","att_type":"jpeg","content_type":"image/jpeg","pid":"6","image_type":"1","key":"","hash":"","bucket":"","status":3,"created_at":"2022-02-08 13:54:50","updated_at":"2022-02-08 13:54:50"},{"att_id":"86","user_id":"1","real_name":"官方商城.png","name":"f06a7bd9-71a9-481c-8f1e-6293bc7a4774-cv0v9.png","att_dir":"2022/02/08/","att_size":"455452","att_type":"png","content_type":"image/png","pid":"6","image_type":"1","key":"","hash":"","bucket":"","status":3,"created_at":"2022-02-08 13:54:50","updated_at":"2022-02-08 13:54:50"},{"att_id":"87","user_id":"1","real_name":"贵人鸟.jpg","name":"e14ce24e-d05d-4d65-86fa-dcaed1fda236-4h6bw.jpeg","att_dir":"2022/02/08/","att_size":"11306","att_type":"jpeg","content_type":"image/jpeg","pid":"6","image_type":"1","key":"","hash":"","bucket":"","status":3,"created_at":"2022-02-08 13:54:50","updated_at":"2022-02-08 13:54:50"},{"att_id":"88","user_id":"1","real_name":"鸿星尔克.jpg","name":"110beeaa-616e-4695-bb86-f4fefaea0503-95ken.jpeg","att_dir":"2022/02/08/","att_size":"11542","att_type":"jpeg","content_type":"image/jpeg","pid":"6","image_type":"1","key":"","hash":"","bucket":"","status":3,"created_at":"2022-02-08 13:54:50","updated_at":"2022-02-08 13:54:50"},{"att_id":"89","user_id":"1","real_name":"回力.jpg","name":"56a57d72-9704-4004-8a72-a455f93ecc3f-teeuz.jpeg","att_dir":"2022/02/08/","att_size":"8972","att_type":"jpeg","content_type":"image/jpeg","pid":"6","image_type":"1","key":"","hash":"","bucket":"","status":3,"created_at":"2022-02-08 13:54:50","updated_at":"2022-02-08 13:54:50"},{"att_id":"90","user_id":"1","real_name":"卡帕.jpg","name":"b985ea57-f440-4372-8266-757ffe98c6df-teeuz.jpeg","att_dir":"2022/02/08/","att_size":"14369","att_type":"jpeg","content_type":"image/jpeg","pid":"6","image_type":"1","key":"","hash":"","bucket":"","status":3,"created_at":"2022-02-08 13:54:50","updated_at":"2022-02-08 13:54:50"},{"att_id":"91","user_id":"1","real_name":"李宁.jpg","name":"e84bf9d0-7a21-43be-93e4-7c5369c5ad4f-teeuz.jpeg","att_dir":"2022/02/08/","att_size":"6047","att_type":"jpeg","content_type":"image/jpeg","pid":"6","image_type":"1","key":"","hash":"","bucket":"","status":3,"created_at":"2022-02-08 13:54:50","updated_at":"2022-02-08 13:54:50"},{"att_id":"92","user_id":"1","real_name":"美津侬.jpg","name":"42c16fd6-894a-4b08-822f-4ec65b24bbde-jmvgt.jpeg","att_dir":"2022/02/08/","att_size":"8221","att_type":"jpeg","content_type":"image/jpeg","pid":"6","image_type":"1","key":"","hash":"","bucket":"","status":3,"created_at":"2022-02-08 13:54:50","updated_at":"2022-02-08 13:54:50"},{"att_id":"93","user_id":"1","real_name":"耐克.jpg","name":"7f46ca65-2fb9-4ceb-9543-c4aa06bfd5c3-qn2yq.jpeg","att_dir":"2022/02/08/","att_size":"14148","att_type":"jpeg","content_type":"image/jpeg","pid":"6","image_type":"1","key":"","hash":"","bucket":"","status":3,"created_at":"2022-02-08 13:54:50","updated_at":"2022-02-08 13:54:50"},{"att_id":"94","user_id":"1","real_name":"匹克.jpg","name":"7f1dad5d-7550-44ff-a4b5-76a4723f51f3-5jt1z.jpeg","att_dir":"2022/02/08/","att_size":"19182","att_type":"jpeg","content_type":"image/jpeg","pid":"6","image_type":"1","key":"","hash":"","bucket":"","status":3,"created_at":"2022-02-08 13:54:51","updated_at":"2022-02-08 13:54:51"},{"att_id":"95","user_id":"1","real_name":"锐步.jpg","name":"537261a7-2a19-46e6-b9f7-d8f108613a2c-eckid.jpeg","att_dir":"2022/02/08/","att_size":"11023","att_type":"jpeg","content_type":"image/jpeg","pid":"6","image_type":"1","key":"","hash":"","bucket":"","status":3,"created_at":"2022-02-08 13:54:51","updated_at":"2022-02-08 13:54:51"},{"att_id":"96","user_id":"1","real_name":"特步.jpg","name":"e163700f-66ba-422f-ad3b-1e97f71ae365-eckid.jpeg","att_dir":"2022/02/08/","att_size":"10458","att_type":"jpeg","content_type":"image/jpeg","pid":"6","image_type":"1","key":"","hash":"","bucket":"","status":3,"created_at":"2022-02-08 13:54:51","updated_at":"2022-02-08 13:54:51"},{"att_id":"97","user_id":"1","real_name":"双星.jpg","name":"bb4e6cf5-f588-4d3a-ac53-e34e7ed6fd5f-eckid.jpeg","att_dir":"2022/02/08/","att_size":"83524","att_type":"jpeg","content_type":"image/jpeg","pid":"6","image_type":"1","key":"","hash":"","bucket":"","status":3,"created_at":"2022-02-08 13:54:51","updated_at":"2022-02-08 13:54:51"},{"att_id":"98","user_id":"1","real_name":"新百伦.jpg","name":"72c509cf-44c2-4f0f-aa75-d2f527929213-v93py.jpeg","att_dir":"2022/02/08/","att_size":"7904","att_type":"jpeg","content_type":"image/jpeg","pid":"6","image_type":"1","key":"","hash":"","bucket":"","status":3,"created_at":"2022-02-08 13:54:51","updated_at":"2022-02-08 13:54:51"},{"att_id":"99","user_id":"1","real_name":"乔丹.jpg","name":"6a43851a-e73b-4b88-9513-82b0fd5b652b-eckid.jpeg","att_dir":"2022/02/08/","att_size":"7654","att_type":"jpeg","content_type":"image/jpeg","pid":"6","image_type":"1","key":"","hash":"","bucket":"","status":3,"created_at":"2022-02-08 13:54:51","updated_at":"2022-02-08 13:54:51"},{"att_id":"100","user_id":"1","real_name":"新品上市.jpg","name":"36e67c36-4d84-44e9-9b34-f8ad24cbb485-7wgjk.jpeg","att_dir":"2022/02/08/","att_size":"30970","att_type":"jpeg","content_type":"image/jpeg","pid":"6","image_type":"1","key":"","hash":"","bucket":"","status":3,"created_at":"2022-02-08 13:54:51","updated_at":"2022-02-08 13:54:51"}],"columns":[{"columnName":"att_id"},{"columnName":"user_id"},{"columnName":"real_name"},{"columnName":"name"},{"columnName":"att_dir"},{"columnName":"att_size"},{"columnName":"att_type"},{"columnName":"content_type"},{"columnName":"pid"},{"columnName":"image_type"},{"columnName":"key"},{"columnName":"hash"},{"columnName":"bucket"},{"columnName":"status"},{"columnName":"created_at"},{"columnName":"updated_at"}]}`

//const sqlSelectResponse2 = `{"msgtype":"response","rows":[{"count":"7"}],"columns":[{"columnName":"count"}]}`
//const sqlSelectResponse2 = `{"msgtype":"response","rows":[],"columns":[{"columnName":"file_id"},{"columnName":"file_former_name"},{"columnName":"file_name"},{"columnName":"file_size"},{"columnName":"user_id"},{"columnName":"file_join_id"},{"columnName":"file_join_type"},{"columnName":"hash"},{"columnName":"bucket"},{"columnName":"status"},{"columnName":"created_at"},{"columnName":"updated_at"},{"columnName":"expire_at"}]}`

func (dc *DatabaseConnections) SqlSelect(req *SqlSelectRequest) *serializer.Response {
	opened := dc.ensureOpened(req.Conid, req.Database)
	response := dc.sendRequest(opened, &containers.EchoMessage{Payload: req.Select, MsgType: "sqlSelect"})
	if response == nil {
		return serializer.Fail("Error executing SQL script")
	}

	if response.Err != nil {
		if existing := findByDatabaseConnection(dc.Opened, req.Conid, req.Database); existing != nil && !existing.Disconnected {
			dc.close(req.Conid, req.Database, false)
		}
		return serializer.Fail(response.Err.Error())
	}

	if response.Payload != nil {
		return serializer.SuccessData(serializer.SUCCESS, map[string]interface{}{
			"msgtype": response.MsgType,
			"rows":    response.Payload,
		})
	}

	unmarshal, err := utility.JsonUnmarshal([]byte(sqlSelectResponse1))
	if err != nil {
		return serializer.Fail(err.Error())
	} else {
		return serializer.SuccessData(serializer.SUCCESS, unmarshal)
	}
}
