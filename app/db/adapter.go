package db

import (
	"fmt"
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

type missingAdapter struct {
	name string
}

func (ma *missingAdapter) Open(ConnectionURL) (Session, error) {
	return nil, fmt.Errorf("upper: Missing adapter %q, did you forget to import it?", ma.name)
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

// LookupAdapter returns a previously registered adapter by name.
func LookupAdapter(name string) Adapter {
	adapterMapMu.RLock()
	defer adapterMapMu.RUnlock()

	if adapter, ok := adapterMap[name]; ok {
		return adapter
	}
	return &missingAdapter{name: name}
}

// Open attempts to stablish a connection with a database.
func Open(adapterName string, settings ConnectionURL) (Session, error) {
	return LookupAdapter(adapterName).Open(settings)
}
