package sideQuests

import "keeper/app/code"

var databaseLastPing code.UnixTime

type DatabaseConnectionHandlers struct {
	Ch chan interface{}
}

func NewDatabaseConnectionHandlers() *DatabaseConnectionHandlers {
	return &DatabaseConnectionHandlers{}
}

func (msg *DatabaseConnectionHandlers) Connect(connection map[string]interface{}) {

}
