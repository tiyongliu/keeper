package bridge

import (
	"context"
	"keeper/app/pkg/serializer"
	"time"
)

type ConnPool struct {
	Conid  string
	Status map[string]string
}

type ServerConnections struct {
	Ctx        context.Context
	Opened     []string
	LastPinged map[string]UnixTime
}

type PingRequest struct {
	Connections []string
}

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

func (sc *ServerConnections) ensureOpened(conid string) {

}

func (sc *ServerConnections) Ping(request *PingRequest) interface{} {
	//_.uniq(connections)

	if request == nil {
		return serializer.Fail(context.Background(), "")
	}

	for _, conid := range request.Connections {
		if l := sc.LastPinged; l != nil {
			last := sc.LastPinged[conid]
			if last != 0 && UnixTime(time.Now().Unix())-last < UnixTime(30*1000) {
				//return Promise.resolve();
			}
			sc.LastPinged[conid] = UnixTime(time.Now().Unix())
		}
	}

	return serializer.SuccessData(Application.ctx, "", map[string]string{"status": "ok"})
}
