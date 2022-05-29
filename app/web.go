package app

import (
	"context"
	"keeper/app/src/keeperContext"
	"keeper/app/src/pkg/logger"
	"keeper/app/src/startup"
	"keeper/app/src/variable"
	"net/http"
	"strconv"
	"time"
)

const localAddr = "0.0.0.0:"

func run(application *variable.SystemApplication) {
	if application == nil || application.Application == nil {
		return
	}

	callRunner(application.Application)
}

func callRunner(webApp *variable.Application) error {
	srv := &http.Server{
		Addr:      localAddr + strconv.Itoa(webApp.Port),
		Handler:   startup.RegisterHttpRoute(),
		TLSConfig: nil,
		//ReadTimeout:       0,
		//ReadHeaderTimeout: 0,
		WriteTimeout:   time.Duration(10) * time.Second,
		IdleTimeout:    time.Duration(10) * time.Second,
		MaxHeaderBytes: 1 << 20,
		//TLSNextProto:      nil,
		//ConnState:         nil,
		//ErrorLog:          nil,
		//BaseContext:       nil,
		//ConnContext:       nil,
	}

	ctx, _ := keeperContext.WithCancelAndSignalHandler()

	go func() {
		<-ctx.Done()

		timeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(timeout); err != nil {
			logger.Fatalf("Server Shutdown Error: %v", err)
		}
		select {
		case <-timeout.Done():
			logger.Error("Shutting down HTTP server timeout")
		default:
			logger.Info("Server exited")
		}
	}()

	return srv.ListenAndServe()
}

/*
func callRunner(webApp *variable.Application) {
	srvHandler := &http.Server{
		Addr:      localAddr + strconv.Itoa(webApp.Port),
		Handler:   startup.RegisterHttpRoute(),
		TLSConfig: nil,
		//ReadTimeout:       0,
		//ReadHeaderTimeout: 0,
		WriteTimeout:   time.Duration(10) * time.Second,
		IdleTimeout:    time.Duration(10) * time.Second,
		MaxHeaderBytes: 1 << 20,
		//TLSNextProto:      nil,
		//ConnState:         nil,
		//ErrorLog:          nil,
		//BaseContext:       nil,
		//ConnContext:       nil,
	}

	ctx, _ := keeperContext.WithCancelAndSignalHandler()

	go func() {
		if err := srvHandler.ListenAndServe(); err != nil {
			log.Printf("Listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := srvHandler.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	log.Println("Server exiting")
}
*/
