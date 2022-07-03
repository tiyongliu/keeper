package modules

type EchoMessage struct {
	Payload interface{} `json:"payload"`
	MsgType string      `json:"msgType"`
	Dialect string      `json:"dialect"`
}
