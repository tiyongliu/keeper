package bridge

import (
	"context"
	"fmt"
	"keeper/app/pkg/logger"
	"sync"
)

var Application *App
var applicationOnce sync.Once

// App struct
type App struct {
	ctx context.Context

	Connections         *Connections
	DatabaseConnections *DatabaseConnections
	ServerConnections   *ServerConnections
}

// NewApp creates a new App application struct
func NewApp() *App {
	applicationOnce.Do(func() {
		Application = &App{
			Connections:         NewConnections(),
			DatabaseConnections: NewDatabaseConnections(),
			ServerConnections:   NewServerConnections(),
		}
	})

	return Application
}

// startup is called when the app starts. The keeperCtx is saved
// so we can call the runtime methods
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
	Application.ctx = ctx
	logger.Infof("Starting up October")
}

// shutdown is called at application termination
func (a *App) Shutdown(ctx context.Context) {
	// Perform your teardown here
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// domReady is called after the front-end dom has been loaded
func (a App) DomReady(ctx context.Context) {
	// Add your action here
}
