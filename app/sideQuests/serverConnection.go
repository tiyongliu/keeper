package sideQuests

type EchoMessage struct {
	Payload interface{}
	MsgType string
}

func SpeakerServerConnection(ch chan *EchoMessage) {
	defer close(ch)

	ch <- &EchoMessage{
		Payload: nil,
		MsgType: "status",
	}

	//TODO 根据params 调用具体的handle下面的方法
}
