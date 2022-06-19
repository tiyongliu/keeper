package bridge

import (
	"context"
	"keeper/app/pkg/serializer"
)

type ServerConnections struct {
	Ctx context.Context
}

const (
	conid = "conid"
)

func NewServerConnections() *ServerConnections {
	return &ServerConnections{}
}

func (sc *ServerConnections) ListDatabases(request map[string]string) interface{} {
	if request["conid"] == "" {
		return serializer.Fail(Application.ctx, "")
	}

	//https://esc.show/article/Golang-GUI-kai-fa-zhi-Webview

	return nil
}
