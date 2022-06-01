package main

import (
	"context"
	"fmt"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"keeper/app/pkg/logger"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The keeperCtx is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	logger.Infof("Starting up October")
}

// shutdown is called at application termination
func (a *App) shutdown(ctx context.Context) {
	// Perform your teardown here
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// domReady is called after the front-end dom has been loaded
func (a App) domReady(ctx context.Context) {
	// Add your action here
}

func (a *App) OpenDirectoryDialog(name string) interface{} {
	selection, err := runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
		Title:         name,
		Message:       "Select a number",
		Buttons:       []string{"one", "two", "three", "four"},
		DefaultButton: "two",
	})

	if err != nil {
		return err.Error()
	}

	return selection
}
