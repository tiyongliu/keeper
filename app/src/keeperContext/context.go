package keeperContext

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"unsafe"
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

type iface struct {
	itab, data uintptr
}

type valueCtx struct {
	context.Context
	key, val interface{}
}

func GetKeyValues(ctx context.Context) map[interface{}]interface{} {
	m := make(map[interface{}]interface{})
	getKeyValue(ctx, m)
	return m
}

func getKeyValue(ctx context.Context, m map[interface{}]interface{}) {
	ictx := *(*iface)(unsafe.Pointer(&ctx))
	if ictx.data == 0 {
		return
	}

	valCtx := (*valueCtx)(unsafe.Pointer(ictx.data))
	if valCtx != nil && valCtx.key != nil && valCtx.val != nil {
		m[valCtx.key] = valCtx.val
	}
	getKeyValue(valCtx.Context, m)
}
