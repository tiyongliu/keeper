package main

import (
	"embed"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
	"keeper/app/bridge"
	"log"
)

//go:embed frontend/dist
var assets embed.FS

func main() {
	//go app.RunApplication()

	// Create an instance of the app structure
	app := bridge.NewApp()
	bridge.Application = app

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "wails",
		Width:  1024,
		Height: 768,
		// MinWidth:          720,
		// MinHeight:         570,
		// MaxWidth:          1280,
		// MaxHeight:         740,
		DisableResize:     false,
		Fullscreen:        false,
		Frameless:         false,
		StartHidden:       false,
		HideWindowOnClose: false,
		RGBA:              &options.RGBA{R: 255, G: 255, B: 255, A: 255},
		Assets:            assets,
		LogLevel:          logger.DEBUG,
		OnStartup:         app.Startup,
		OnDomReady:        app.DomReady,
		OnShutdown:        app.Shutdown,
		OnBeforeClose:     app.OnBeforeClose,
		Bind: []interface{}{
			app,
			app.Connections,
			app.DatabaseConnections,
			app.ServerConnections,
			app.Plugins,
		},
		// Windows platform specific options
		Windows: &windows.Options{
			WebviewIsTransparent: false,
			WindowIsTranslucent:  false,
			DisableWindowIcon:    false,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
}

/*
#### A Lodash-style Go library based on Go 1.18+ Generics (map, filter, contains, find...)
[go的lodash](https://github.com/samber/lo)

[gui-插件](https://github.com/go-graphics/go-gui-projects.git)

[go windows gui api](https://github.com/rodrigocfd/windigo)

[着色器](https://github.com/gen2brain/iup-go/tree/main/examples/getcolor)

[小而巧的跨平台原生GUI库](https://z-kit.cc/)

[Golang GUI开发之Webview](https://esc.show/article/Golang-GUI-kai-fa-zhi-Webview)

[QT GIF 编辑器](https://gitee.com/wingsummer/wing-gif-editor)

wails update -pre
*/
