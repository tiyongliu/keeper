package bridge

import (
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"keeper/app/db/adapter/mongo"
	"keeper/app/db/adapter/mysql"
	"keeper/app/pkg/serializer"
)

type Plugins struct {
}

func NewPlugins() *Plugins {
	return &Plugins{}
}

func (p *Plugins) Installed() *serializer.Response {
	return serializer.SuccessData(serializer.SUCCESS, []map[string]string{
		{"name": mongo.Adapter},
		{"name": mysql.Adapter},
	})
}

type ScriptRequest struct {
	PackageName string `json:"packageName"`
}

func (p *Plugins) Script(req *ScriptRequest) *serializer.Response {
	module := make(chan interface{}, 2)
	runtime.EventsEmit(Application.ctx, "pullEventPluginsScript", req.PackageName)
	runtime.EventsOnce(Application.ctx, "loadPlugins", func(resp ...interface{}) {
		if resp != nil {
			module <- resp[0]
		}
	})
	defer close(module)
	return serializer.SuccessData(serializer.SUCCESS, <-module)
}
