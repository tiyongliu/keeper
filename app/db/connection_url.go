package db

// ConnectionURL represents a data source name (DSN).
type ConnectionURL interface {
	// String returns the connection string that is going to be passed to the
	// adapter.
	String() string
}
