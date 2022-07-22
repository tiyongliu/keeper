package utility

import (
	"context"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func EmitChanged(ctx context.Context, key string) {
	runtime.EventsEmit(ctx, "clean-cache", key)
	runtime.EventsEmit(ctx, key)
}
