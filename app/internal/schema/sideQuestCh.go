package schema

import (
	"keeper/app/db/standard/modules"
	"keeper/app/utility"
)

type EchoMessage struct {
	Payload interface{}
	MsgType string
	Dialect string
	Err     error
}

type OpenedStatus struct {
	Name    string `json:"name"`
	Message string `json:"message"`
	Counter int    `json:"counter"`
}

type OpenedServerConnection struct {
	Conid        string                 `json:"conid"`
	Status       *OpenedStatus          `json:"status"`
	Databases    interface{}            `json:"databases"`
	Connection   map[string]interface{} `json:"connection"`
	Disconnected bool                   `json:"disconnected"`
	Version      *modules.Version       `json:"version"`
}

type OpenedDatabaseConnection struct {
	Conid         string                 `json:"conid"`
	Status        *OpenedStatus          `json:"status"`
	Database      string                 `json:"database"`
	Connection    map[string]interface{} `json:"connection"`
	ServerVersion *modules.Version       `json:"version"`
	Structure     map[string]interface{} `json:"structure"`
	AnalysedTime  utility.UnixTime       `json:"analysedTime"`
	Disconnected  bool                   `json:"disconnected"`
}

type DatabaseConnectionClosed struct {
	Structure    map[string]interface{} `json:"structure"`
	AnalysedTime utility.UnixTime       `json:"analysedTime"`
	Status       *OpenedStatus          `json:"status"`
}
