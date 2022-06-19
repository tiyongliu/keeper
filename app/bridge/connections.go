package bridge

import (
	"context"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/sqweek/dialog"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"keeper/app/modules"
	"keeper/app/pkg/serializer"
	"keeper/app/pkg/standard"
	plugin_mondb "keeper/app/plugins/plugin-mondb"
	plugin_mysql "keeper/app/plugins/plugin-mysql"
	"keeper/app/tools"
	"keeper/app/utility"
	"path"
	"sync"
)

var ConnectionsBridge *Connections

var bridgeOnce sync.Once

var JsonLinesDatabase *utility.JsonLinesDatabase

type Connections struct {
	Ctx context.Context
}

const (
	mysql_alias = "mysql"
	mongo_alias = "mongo"
)

func init() {
	JsonLinesDatabase = utility.NewJsonLinesDatabase(path.Join(tools.DataDirCore(), "connections.jsonl"))
}

func NewConnections() *Connections {
	bridgeOnce.Do(func() {
		ConnectionsBridge = &Connections{}
	})

	return ConnectionsBridge
}

func (conn *Connections) Test(connection map[string]interface{}) interface{} {
	if connection["engine"].(string) == mysql_alias {
		simpleSettingMysql := &modules.SimpleSettingMysql{}
		err := mapstructure.Decode(connection, simpleSettingMysql)
		if err != nil {
			return err.Error()
		}

		pool, err := plugin_mysql.NewSimpleMysqlPool(simpleSettingMysql)

		if err != nil {
			selection, _ := runtime.MessageDialog(Application.ctx, runtime.MessageDialogOptions{
				Type:          runtime.ErrorDialog,
				Title:         "测试失败",
				Message:       err.Error(),
				Buttons:       []string{"确认"},
				DefaultButton: "确认",
			})

			return selection
		}

		defer pool.Close()

		getVersion, err := pool.GetVersion()
		if err != nil {
			selection, _ := runtime.MessageDialog(Application.ctx, runtime.MessageDialogOptions{
				Type:          runtime.ErrorDialog,
				Title:         "测试失败",
				Message:       err.Error(),
				Buttons:       []string{"确认"},
				DefaultButton: "确认",
			})
			return selection
		}

		driver := getVersion.(*standard.VersionMsg)
		selection, _ := runtime.MessageDialog(Application.ctx, runtime.MessageDialogOptions{
			Title:         "连接成功",
			Message:       "Connected" + fmt.Sprintf(": %s", driver.VersionText),
			Buttons:       []string{"确认"},
			DefaultButton: "确认",
		})
		return selection

	} else if connection["engine"].(string) == mongo_alias {
		pool, err := plugin_mondb.NewSimpleMongoDBPool(&modules.SimpleSettingMongoDB{
			Host: connection["host"].(string),
			Port: connection["port"].(string),
		})

		defer pool.Close()

		if err != nil {
			runtime.MessageDialog(Application.ctx, runtime.MessageDialogOptions{
				Type:          runtime.ErrorDialog,
				Title:         "测试失败",
				Message:       err.Error(),
				Buttons:       []string{"确认"},
				DefaultButton: "确认",
			})
			return err.Error()
		}

		getVersion, err := pool.GetVersion()
		if err != nil {
			runtime.MessageDialog(Application.ctx, runtime.MessageDialogOptions{
				Type:          runtime.ErrorDialog,
				Title:         "测试失败",
				Message:       err.Error(),
				Buttons:       []string{"确认"},
				DefaultButton: "确认",
			})
			return err.Error()
		}

		driver := getVersion.(*standard.VersionMsg)
		selection, _ := runtime.MessageDialog(Application.ctx, runtime.MessageDialogOptions{
			Title:         "连接成功",
			Message:       "Connected" + fmt.Sprintf(": %s", driver.VersionText),
			Buttons:       []string{"确认"},
			DefaultButton: "确认",
		})
		return selection
	}

	return nil
}

func (conn *Connections) Save(connection map[string]string) interface{} {
	encrypted := utility.EncryptConnection(connection)
	//验证obj的唯一性，除去key字段，所有key对应的值都要一致。
	unknownMap := tools.TransformUnknownMap(encrypted)
	if exists := tools.UnknownMapSome(JsonLinesDatabase.Find(), unknownMap); exists {
		runtime.MessageDialog(Application.ctx, runtime.MessageDialogOptions{
			Type:          runtime.ErrorDialog,
			Title:         "错误",
			Message:       "Connection with same connection name already exists in the project.",
			Buttons:       []string{"确认"},
			DefaultButton: "确认",
		})
		return serializer.Fail(Application.ctx, "")
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
			Title:         "连接失败",
			Message:       err.Error(),
			Buttons:       []string{"确认"},
			DefaultButton: "确认",
		})

		return serializer.Fail(Application.ctx, "")
	}

	return serializer.SuccessData(Application.ctx, "", res)
}

func (conn *Connections) List() interface{} {
	find := JsonLinesDatabase.Find()
	return serializer.SuccessData(Application.ctx, "", find)
}

func (conn *Connections) Get(conid map[string]string) interface{} {
	return serializer.SuccessData(Application.ctx, "", nil)
}

func (conn *Connections) getCore(conid string, mask bool) {
	if conid == "" {
		return
	}
}

func (conn *Connections) Delete(connection map[string]string) interface{} {
	ok := dialog.Message("你确认要删除\"%s\"吗?", connection["name"]).Title("确认删除").YesNo()
	if !ok {
		return serializer.Fail(Application.ctx, fmt.Sprintf("%v", ok))
	}

	uuid, ok := connection["_id"]
	if ok && uuid != "" {
		res, err := JsonLinesDatabase.Remove(uuid)
		if err != nil {
			runtime.MessageDialog(Application.ctx, runtime.MessageDialogOptions{
				Type:          runtime.ErrorDialog,
				Title:         "删除失败",
				Message:       err.Error(),
				Buttons:       []string{"确认"},
				DefaultButton: "确认",
			})

			return serializer.Fail(Application.ctx, err.Error())
		}

		runtime.EventsEmit(Application.ctx, "connection-list-changed", res)
		return serializer.SuccessData(Application.ctx, "", res)
	}

	return serializer.Fail(Application.ctx, "参数错误")
}
