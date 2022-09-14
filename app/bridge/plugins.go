package bridge

import (
	"keeper/app/pkg/serializer"
	"keeper/app/pkg/standard"
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
		{"name": standard.MONGOALIAS},
		{"name": standard.MYSQLALIAS},
	})
}

type ScriptRequest struct {
	PackageName string `json:"packageName"`
}

func (p *Plugins) Script(req *ScriptRequest) *serializer.Response {
	switch req.PackageName {
	case standard.MONGOALIAS:
		return serializer.SuccessData(serializer.SUCCESS, map[string]interface{}{"drivers": MongoDBFrontend.Driver()})
	case standard.MYSQLALIAS:
		return serializer.SuccessData(serializer.SUCCESS, map[string]interface{}{"drivers": MysqlFrontend.Driver()})
	default:
		return serializer.SuccessData(serializer.SUCCESS, nil)
	}
}
