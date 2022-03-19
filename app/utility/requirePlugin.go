package utility

import (
	"errors"
	keeper_plugin_mongo "keeper/plugins"
)

var loadedPlugins map[string]*keeper_plugin_mongo.DriversPlugin

func requirePlugin(packageName string, requiredPlugin *keeper_plugin_mongo.DriversPlugin) (*keeper_plugin_mongo.DriversPlugin, error) {
	if packageName == "" {
		return nil, errors.New("Missing packageName in plugin")
	}

	if loadedPlugins[packageName] != nil {
		return loadedPlugins[packageName], nil
	}

	if requiredPlugin == nil {

	}

	loadedPlugins[packageName] = keeper_plugin_mongo.RequiredPlugin

	return requiredPlugin, nil
}
