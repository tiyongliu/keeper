package db

import (
	"sync"
)

var (
	adapterMap   = make(map[string]Adapter)
	adapterMapMu sync.RWMutex
)

// Adapter interface defines an adapter
type Adapter interface {
	Open(ConnectionURL) (Session, error)
}

// RegisterAdapter registers a generic database adapter.
func RegisterAdapter(name string, adapter Adapter) {
	adapterMapMu.Lock()
	defer adapterMapMu.Unlock()

	if name == "" {
		panic(`Missing adapter name`)
	}
	if _, ok := adapterMap[name]; ok {
		panic(`db.RegisterAdapter() called twice for adapter: ` + name)
	}
	adapterMap[name] = adapter
}
