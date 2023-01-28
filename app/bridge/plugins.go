package bridge

import (
	"keeper/app/db/adapter/mongo"
	"keeper/app/db/adapter/mysql"
	"keeper/app/pkg/serializer"
	MongoDBFrontend "keeper/app/plugins/pluginMongdb/frontend"
	MysqlFrontend "keeper/app/plugins/pluginMysql/frontend"
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
	switch req.PackageName {
	case mongo.Adapter:
		return serializer.SuccessData(serializer.SUCCESS, map[string]interface{}{"drivers": MongoDBFrontend.Driver()})
	case mysql.Adapter:
		return serializer.SuccessData(serializer.SUCCESS, map[string]interface{}{"drivers": MysqlFrontend.Driver()})
	default:
		return serializer.SuccessData(serializer.SUCCESS, nil)
	}
}
