package bridge

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"keeper/app/internal"
	"keeper/app/pkg/serializer"
	"keeper/app/pkg/standard"
	"keeper/app/plugins/modules"
	"keeper/app/plugins/pluginMongdb"
	"keeper/app/plugins/pluginMysql"
	"keeper/app/utility"
	"path"
)

var JsonLinesDatabase *utility.JsonLinesDatabase

type Connections struct {
}

func init() {
	dir := utility.DataDir()
	JsonLinesDatabase = utility.NewJsonLinesDatabase(path.Join(dir, "connections.jsonl"))
}

func NewConnections() *Connections {
	return &Connections{}
}

func getCore(conid string, mask bool) map[string]interface{} {
	if conid == "" {
		return nil
	}
	return JsonLinesDatabase.Get(conid)
}

const (
	testTitleFailed   = "Test failed"
	testTitleSuccess  = "Test success"
	connectionFailed  = "Connection failed"
	connectionSuccess = "Connection success"
	deleteFailed      = "Delete failed"
	oK                = "OK"
)

func (conn *Connections) Test(connection map[string]interface{}) *serializer.Response {
	switch connection["engine"].(string) {
	case standard.MYSQLALIAS:
		simpleSettingMysql := &modules.SimpleSettingMysql{}
		err := mapstructure.Decode(connection, simpleSettingMysql)
		if err != nil {
			runtime.MessageDialog(Application.ctx, runtime.MessageDialogOptions{
				Type:          runtime.ErrorDialog,
				Title:         "Error",
				Message:       err.Error(),
				Buttons:       []string{oK},
				DefaultButton: oK,
			})
			return serializer.Fail(err.Error())
		}

		pool, err := pluginMysql.NewSimpleMysqlPool(simpleSettingMysql)
		if err != nil {
			runtime.MessageDialog(Application.ctx, runtime.MessageDialogOptions{
				Type:          runtime.ErrorDialog,
				Title:         testTitleFailed,
				Message:       err.Error(),
				Buttons:       []string{oK},
				DefaultButton: oK,
			})
			return serializer.Fail(err.Error())
		} else {
			defer pool.Close()
			driver, err := pool.GetVersion()
			if err != nil {
				runtime.MessageDialog(Application.ctx, runtime.MessageDialogOptions{
					Type:          runtime.ErrorDialog,
					Title:         testTitleFailed,
					Message:       err.Error(),
					Buttons:       []string{oK},
					DefaultButton: oK,
				})
				return serializer.Fail(err.Error())
			}
			runtime.MessageDialog(Application.ctx, runtime.MessageDialogOptions{
				Title:         testTitleSuccess,
				Message:       "Connected" + fmt.Sprintf(": %s", driver.VersionText),
				Buttons:       []string{oK},
				DefaultButton: oK,
			})
			return serializer.SuccessData(serializer.SUCCESS, nil)
		}
	case standard.MONGOALIAS:
		pool, err := pluginMongdb.NewSimpleMongoDBPool(&modules.SimpleSettingMongoDB{
			Host: connection["host"].(string),
			Port: connection["port"].(string),
		})
		if err != nil {
			runtime.MessageDialog(Application.ctx, runtime.MessageDialogOptions{
				Type:          runtime.ErrorDialog,
				Title:         testTitleFailed,
				Message:       err.Error(),
				Buttons:       []string{oK},
				DefaultButton: oK,
			})
			return serializer.Fail(err.Error())
		} else {
			defer pool.Close()
			driver, err := pool.GetVersion()
			if err != nil {
				runtime.MessageDialog(Application.ctx, runtime.MessageDialogOptions{
					Type:          runtime.ErrorDialog,
					Title:         testTitleFailed,
					Message:       err.Error(),
					Buttons:       []string{oK},
					DefaultButton: oK,
				})
				return serializer.Fail(err.Error())
			}

			runtime.MessageDialog(Application.ctx, runtime.MessageDialogOptions{
				Title:         connectionSuccess,
				Message:       "Connected" + fmt.Sprintf(": %s", driver.VersionText),
				Buttons:       []string{oK},
				DefaultButton: oK,
			})

			return serializer.SuccessData(serializer.SUCCESS, nil)
		}
	default:
		return serializer.Fail(serializer.ParamsErr)
	}
}

func (conn *Connections) Save(connection map[string]string) *serializer.Response {
	encrypted := internal.EncryptConnection(connection)
	//验证obj的唯一性，除去key字段，所有key对应的值都要一致。
	unknownMap := utility.TransformUnknownMap(encrypted)
	if exists := utility.UnknownMapSome(JsonLinesDatabase.Find(), unknownMap); exists {
		runtime.MessageDialog(Application.ctx, runtime.MessageDialogOptions{
			Type:          runtime.ErrorDialog,
			Title:         connectionFailed,
			Message:       "Connection with same connection name already exists in the project.",
			Buttons:       []string{oK},
			DefaultButton: oK,
		})
		return serializer.Fail(serializer.ParamsErr)
	}

	uuid, ok := connection["_id"]
	var res map[string]interface{}
	var err error

	if ok && uuid != "" {
		res, err = JsonLinesDatabase.Update(unknownMap)
	} else {
		res, err = JsonLinesDatabase.Insert(unknownMap)
	}

	if err != nil {
		runtime.MessageDialog(Application.ctx, runtime.MessageDialogOptions{
			Type:          runtime.ErrorDialog,
			Title:         testTitleFailed,
			Message:       err.Error(),
			Buttons:       []string{oK},
			DefaultButton: oK,
		})

		return serializer.Fail(serializer.ParamsErr)
	}
	utility.EmitChanged(Application.ctx, "connection-list-changed")
	return serializer.SuccessData(serializer.SUCCESS, res)
}

func (conn *Connections) List() *serializer.Response {
	find := JsonLinesDatabase.Find()
	return serializer.SuccessData(serializer.SUCCESS, find)
}

type GetConnectionsRequest struct {
	Conid string `json:"conid"`
}

func (conn *Connections) Get(req *GetConnectionsRequest) *serializer.Response {
	return serializer.SuccessData(serializer.SUCCESS, getCore(req.Conid, true))
}

func (conn *Connections) Delete(connection map[string]string) *serializer.Response {
	//ok := dialog.Message("你确认要删除\"%s\"吗?", connection["name"]).Title("确认删除").YesNo()
	//if !ok {
	//	return serializer.Fail(Application.ctx, fmt.Sprintf("%v", ok))
	//}

	uuid, ok := connection["_id"]
	if ok && uuid != "" {
		res, err := JsonLinesDatabase.Remove(uuid)
		if err != nil {
			runtime.MessageDialog(Application.ctx, runtime.MessageDialogOptions{
				Type:          runtime.ErrorDialog,
				Title:         deleteFailed,
				Message:       err.Error(),
				Buttons:       []string{oK},
				DefaultButton: oK,
			})

			return serializer.Fail(err.Error())
		}
		utility.EmitChanged(Application.ctx, "connection-list-changed")
		return serializer.SuccessData(serializer.SUCCESS, res)
	}

	return serializer.Fail(serializer.ParamsErr)
}
