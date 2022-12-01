package db

import "keeper/app/db/standard"

// Session is an interface that defines methods for database adapters.
type Session interface {
	Settings

	standard.SQL
}
