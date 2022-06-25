package driver

import (
	"keeper/app/bridge"
	"keeper/app/tools"
	"time"
)

var lastStatus string
var lastPing bridge.UnixTime

type MysqlMessage struct {
}

type StatusMessage struct {
	Name string
}

func NewMysqlMessage() MessageDriverHandlers {
	return &MysqlMessage{}
}

func (msg *MysqlMessage) Connect(Connection map[string]interface{}) {
	setStatusName("pending")
	lastPing = bridge.UnixTime(time.Now().Unix())

	//TODO request to dbEngineDriver
}

func (msg *MysqlMessage) Ping() bridge.UnixTime {
	return bridge.UnixTime(time.Now().Unix())
}

func (msg *MysqlMessage) CreateDatabase() {

}

func setStatusName(name string) {
	setStatus(&StatusMessage{name})

}

func setStatus(status *StatusMessage) {
	statusString := tools.ToJsonStr(status)
	if lastStatus != statusString {
		//TODO send
		lastStatus = statusString
	}
}
