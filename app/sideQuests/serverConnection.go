package sideQuests

import (
	"keeper/app/handler"
	"keeper/app/modules"
)

func SpeakerServerConnection(ch chan *modules.EchoMessage, connection map[string]interface{}) {
	defer close(ch)

	handler.NewMessageDriverHandlers(ch).Connect(connection)

	//ch <- &EchoMessage{
	//	Payload: nil,
	//	MsgType: "status",
	//}

	//TODO 根据params 调用具体的handle下面的方法
}
