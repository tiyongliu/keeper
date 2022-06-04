package bridge

import (
	"context"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"keeper/app/modules"
	"keeper/app/pkg/standard"
	plugin_mondb "keeper/app/plugins/plugin-mondb"
	plugin_mysql "keeper/app/plugins/plugin-mysql"
	"keeper/app/utility"
	"sync"
)

var ConnectionsBridge *Connections

var bridgeOnce sync.Once

type Connections struct {
	Ctx context.Context
}

const (
	mysql_alias = "mysql"
	mongo_alias = "mongo"
)

func NewConnectProcess() *Connections {
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
	fmt.Printf("encrypted: %s\n", encrypted)

	_id := connection["_id"]

	fmt.Println("_id: ", _id)
	fmt.Printf("_id 134:  %v \n", _id)

	return connection
}
