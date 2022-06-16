package bridge

import "context"

type ServerConnections struct {
	Ctx context.Context
}

func NewServerConnections() *ServerConnections {
	return &ServerConnections{}
}
