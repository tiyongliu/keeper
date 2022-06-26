package utility

import (
	"keeper/app/pkg/logger"
	"reflect"
	"strings"
)

func RequireEngineDriver(connection interface{}) {
	var isEngine bool
	var engine string
	//connection
	_Kind := reflect.TypeOf(connection).Kind()
	if _Kind == reflect.String {
		engine = connection.(string)
		isEngine = true
	} else if _Kind == reflect.Map {
		engine = connection.(map[string]string)["engine"]
		isEngine = true
	}
	if engine == "" && !isEngine {
		logger.Fatalf("Could not get driver from connection \n")
	}

	if strings.Contains(engine, "@") {
		split := strings.Split(engine, "@")
		//shortName := split[0]
		packageName := split[1]
		plugin, err := requirePlugin(packageName, nil)
		if err != nil {
			return
		}

		if plugin != nil && len(plugin.Drivers) > 0 {
			for _, x := range plugin.Drivers {
				if x.Engine == engine {

				}
			}
		}

	}

	logger.Fatalf("Could not find engine driver %s", engine)
}
