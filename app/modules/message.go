package modules

type EchoMessage struct {
	Payload interface{} `json:"payload"`
	MsgType string      `json:"msgType"`
}

type DriverPayload struct {
	Name        string      `json:"name"`
	StandardRes interface{} `json:"standardRes"`
}
