package modules

type EchoMessage struct {
	Payload interface{} `json:"payload"`
	MsgType string      `json:"msgType"`
	Dialect string      `json:"dialect"`
	Conid   string      `json:"conid"`
}

type OpenedStatus struct {
	Name    string `json:"name"`
	Message string `json:"message"`
}
