package main

import (
	"embed"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"keeper/app/bridge"
	"log"
)

//go:embed frontend/dist
var assets embed.FS

func main() {
	app := bridge.NewApp()
	bridge.Application = app

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "keeper",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		LogLevel:         logger.DEBUG,
		OnStartup:        app.Startup,
		OnDomReady:       app.DomReady,
		OnShutdown:       app.Shutdown,
		OnBeforeClose:    app.OnBeforeClose,
		Bind: []interface{}{
			app,
			app.Connections,
			app.DatabaseConnections,
			app.ServerConnections,
			app.Plugins,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
}
