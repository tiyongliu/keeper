package services

import (
	"fmt"
	"keeper/app/modules"
	plugin_mondb "keeper/app/plugins/plugin-mondb"
)

func ProcessMessage(connection *modules.MysqlConnection) {
	if connection == nil {
		return
	}

	//connectUtility
	//getVersion
}

func GetVersion(connection map[string]string) {
	version, err := plugin_mondb.NewMongoDB().GetVersion()
	fmt.Println(version)
	fmt.Println(err)
}
