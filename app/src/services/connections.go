package services

import (
	"dbbox/app/src/modules"
	plugin_mondb "dbbox/app/src/plugins/plugin-mondb"
	"fmt"
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
