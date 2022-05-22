package main

import (
	"embed"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed packages/web/dist
var assets embed.FS

func web() {
	engine := gin.Default()
	engine.GET("/", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"msg": 1,
		})
	})
	engine.Run(":34115")
}
func main() {

	go func() {
		web()
	}()

	// Create an instance of the app structure
	app := NewApp()

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
		OnStartup:         app.startup,
		OnDomReady:        app.domReady,
		OnShutdown:        app.shutdown,
		Bind: []interface{}{
			app,
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
