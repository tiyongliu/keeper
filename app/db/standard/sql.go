package standard

import "keeper/app/db/standard/modules"

// Standard
type SQL interface {
	Dialect() string
	Ping() error
	Version() (*modules.Version, error)
	Close() error
	ListDatabases() (interface{}, error)
}
