package containers

import (
	"keeper/app/pkg/standard"
)

type EchoMessage struct {
	Payload interface{} `json:"payload"`
	MsgType string      `json:"msgType"`
	Dialect string      `json:"dialect"`
}

type OpenedStatus struct {
	Name    string `json:"name"`
	Message string `json:"message"`
}

type OpenedData struct {
	Conid        string                 `json:"conid"`
	Status       *OpenedStatus          `json:"status"`
	Databases    interface{}            `json:"databases"`
	Connection   map[string]interface{} `json:"connection"`
	Disconnected bool                   `json:"disconnected"`
	Version      *standard.VersionMsg   `json:"version"`
}
