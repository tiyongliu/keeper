package driver

import "keeper/app/bridge"

type MessageDriverHandlers interface {
	Connect(Connection map[string]interface{})
	Ping() bridge.UnixTime
	CreateDatabase()
}
