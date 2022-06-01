package main

/*

var Application *App
var applicationOnce sync.Once

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	applicationOnce.Do(func() {
		Application = &App{}
	})
	return Application
}

// startup is called when the app starts. The keeperCtx is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	Application.ctx = ctx
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
		Type:          runtime.InfoDialog,
		Title:         "连接成功",
		Message:       "连接成功",
		Buttons:       []string{"确认"},
		DefaultButton: "确认",
	})
	//selection, err := runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
	//	Type:          runtime.ErrorDialog,
	//	Title:         "测试失败",
	//	Message:       "1045 - Access denied for user 'roo'@'localhost' (using password: YES)",
	//	Buttons:       []string{"确认"},
	//	DefaultButton: "确认",
	//})

	if err != nil {
		return err.Error()
	}

	return selection
}

*/
