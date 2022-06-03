package bridge

import (
	"context"
	"fmt"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"keeper/app/modules"
	"keeper/app/pkg/standard"
	plugin_mondb "keeper/app/plugins/plugin-mondb"
	plugin_mysql "keeper/app/plugins/plugin-mysql"
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

func (conn *Connections) Test(req map[string]interface{}) interface{} {
	if req["engine"].(string) == mysql_alias {
		pool, err := plugin_mysql.NewSimpleMysqlPool(&modules.SimpleSettingMysql{
			Host:     req["server"].(string),
			Username: req["user"].(string),
			Password: req["password"].(string),
			Port:     req["port"].(string),
		})

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

	} else if req["engine"].(string) == mongo_alias {
		pool, err := plugin_mondb.NewSimpleMongoDBPool(&modules.SimpleSettingMongoDB{
			Host: req["server"].(string),
			Port: req["port"].(string),
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

func (conn *Connections) Save() {

}
