package cmd

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"keeper/app/src/api"
	"keeper/app/src/pkg/logger"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	addReqInfo = "additionalReqInfo"
)

type ApiServer interface {
	Run(context.Context, chan struct{}) error
	RegisterApi(*api.ApiDefinition) *RegistResult
	RegisterMiddleWare(func(*gin.Context))
}

type apiServerOptions struct {
	name string
	addr string
}

type apiServer struct {
	sync.Mutex
	options *apiServerOptions
	apiDefs map[string]*api.ApiDefinition
	engine  *gin.Engine
}

type ApiServerOption func(*apiServerOptions)

func WithName(name string) ApiServerOption {
	return func(options *apiServerOptions) {
		options.name = name
	}
}

func WithAddr(addr string) ApiServerOption {
	return func(options *apiServerOptions) {
		options.addr = addr
	}
}

func newApiServer(options ...ApiServerOption) ApiServer {
	o := &apiServerOptions{}
	for _, option := range options {
		option(o)
	}

	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.GET("/ping", handle("/ping", &api.ApiHandler{
		Handler: api.Ping,
		Method:  api.GET,
		Pattern: "ping",
	}, &api.ApiDefinition{}))

	as := &apiServer{
		options: o,
		engine:  engine,
	}

	return as
}

func (as *apiServer) Run(ctx context.Context, shuttedDown chan struct{}) error {
	logger.Infof("Start ApiServer[%s] on [%s]", as.options.name, as.options.addr)
	srv := http.Server{
		Addr:    as.options.addr,
		Handler: as.engine,
	}

	go func() {
		<-ctx.Done()
		defer func() {
			shuttedDown <- struct{}{}
		}()
		logger.Info("Shutting down HTTP server ...")

		timeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(timeout); err != nil {
			logger.Fatalf("Server Shutdown Error: %v", err)
		}
		select {
		case <-timeout.Done():
			logger.Warnf("Shutting down HTTP server timeout")
		default:
			logger.Info("Server exited")
		}
	}()

	return srv.ListenAndServe()
}

type RegistResult struct {
	Err        error
	Name       string
	PathPrefix string
}

func (r *RegistResult) setErr(err error) *RegistResult {
	r.Err = err
	return r
}

func (as *apiServer) RegisterApi(apiDef *api.ApiDefinition) *RegistResult {
	as.Lock()
	defer as.Unlock()
	result := &RegistResult{
		Name:       apiDef.Name,
		PathPrefix: apiDef.PathPrefix,
	}

	if len(apiDef.Name) == 0 {
		return result.setErr(errors.New("ApiDefinition name is null"))
	}

	if len(apiDef.PathPrefix) == 0 {
		return result.setErr(fmt.Errorf("No pathPrefix is specified for %s", apiDef.Name))
	}

	_, ok := as.apiDefs[apiDef.Name]

	if ok {
		return result.setErr(fmt.Errorf("Duplicate ApiDefinition found: %s", apiDef.Name))
	}

	return result.setErr(as.registerApiHandlers(apiDef))
}

func (as *apiServer) RegisterMiddleWare(handler func(*gin.Context)) {
	as.engine.RouterGroup.Use(handler)
}

func (as *apiServer) registerApiHandlers(apiDef *api.ApiDefinition) error {
	for _, handler := range apiDef.Handlers {
		var err error
		switch handler.Method {
		case api.GET:
			err = as.registerGet(apiDef, handler)
		case api.POST:
			err = as.registerPost(apiDef, handler)
		case api.PUT:
			err = as.registerPut(apiDef, handler)
		case api.DELETE:
			err = as.registerDelete(apiDef, handler)
		default:
			err = fmt.Errorf("Unknown handler method[%s -> %s]", handler.Method.String(), handler.Pattern)
		}

		if err != nil {
			return err
		}
	}
	return nil
}

func methodNotMatchErr(expected api.Method, apiDef *api.ApiDefinition, handler *api.ApiHandler) error {
	return fmt.Errorf("%s -> %s wanted method[%s], got %s", apiDef.Name, handler.Pattern, expected.String(), handler.Method.String())
}

func buildPath(apiDef *api.ApiDefinition, handler *api.ApiHandler) string {
	prefix := strings.ToLower(apiDef.PathPrefix)
	path := strings.ToLower(handler.Pattern)

	if !strings.HasPrefix(prefix, "/") {
		prefix = "/" + prefix
	}

	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	return prefix + path
}

func handle(path string, handler *api.ApiHandler, apiDef *api.ApiDefinition) func(*gin.Context) {
	return func(c *gin.Context) {

	}
}

func (as *apiServer) registerGet(apiDef *api.ApiDefinition, handler *api.ApiHandler) error {
	if handler.Method != api.GET {
		return methodNotMatchErr(api.GET, apiDef, handler)
	}

	fullPath := buildPath(apiDef, handler)
	logger.Infof("Register API handler: [%s %s]", handler.Method.String(), fullPath)

	as.engine.GET(fullPath, handle(fullPath, handler, apiDef))
	return nil
}

func (as *apiServer) registerPost(apiDef *api.ApiDefinition, handler *api.ApiHandler) error {
	if handler.Method != api.POST {
		return methodNotMatchErr(api.POST, apiDef, handler)
	}
	fullPath := buildPath(apiDef, handler)
	logger.Infof("Register API handler: [%s %s]", handler.Method.String(), fullPath)
	as.engine.POST(fullPath, handle(fullPath, handler, apiDef))

	return nil
}

func (as *apiServer) registerPut(apiDef *api.ApiDefinition, handler *api.ApiHandler) error {
	if handler.Method != api.PUT {
		return methodNotMatchErr(api.PUT, apiDef, handler)
	}
	fullPath := buildPath(apiDef, handler)
	logger.Infof("Register API handler: [%s %s]", handler.Method.String(), fullPath)
	as.engine.POST(fullPath, handle(fullPath, handler, apiDef))

	return nil
}

func (as *apiServer) registerDelete(apiDef *api.ApiDefinition, handler *api.ApiHandler) error {
	if handler.Method != api.DELETE {
		return methodNotMatchErr(api.DELETE, apiDef, handler)
	}

	fullPath := buildPath(apiDef, handler)
	logger.Infof("Register API handler: [%s %s]", handler.Method.String(), fullPath)
	as.engine.DELETE(fullPath, handle(fullPath, handler, apiDef))

	return nil
}
