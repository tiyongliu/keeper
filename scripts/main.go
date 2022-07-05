package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func WithCancelAndSignalHandler() (context.Context, context.CancelFunc) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		<-signals
		cancel()
	}()

	return ctx, cancel
}

func main() {
	ctx, _ := WithCancelAndSignalHandler()
	ticker := time.NewTicker(1 * time.Second)

	go func() {
		for range ticker.C {
			nowTime := time.Now().Unix()
			fmt.Println("Tick at", nowTime)
			//ticker.Stop()
		}
	}()

	<-ctx.Done()
}
