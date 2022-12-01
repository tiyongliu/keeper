package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"keeper/app/db"
	"sync"
	"time"
)

// Adapter holds the name of the mongodb adapter.
const Adapter = `mongo`

var connTimeout = time.Second * 5

// Source represents a MongoDB database.
type Source struct {
	db.Settings
	ctx     context.Context
	name    string
	connURL db.ConnectionURL
	//session  *mongo.Session
	//database *mongo.Database
	client        *mongo.Client
	collectionsMu sync.Mutex
}

type mongoAdapter struct {
}

func (mongoAdapter) Open(dsn db.ConnectionURL) (db.Session, error) {
	return Open(dsn)
}

func init() {
	db.RegisterAdapter(Adapter, db.Adapter(&mongoAdapter{}))
}

// Open stablishes a new connection to a SQL server.
func Open(settings db.ConnectionURL) (db.Session, error) {
	d := &Source{Settings: db.NewSettings(), ctx: context.Background()}
	if err := d.Open(settings); err != nil {
		return nil, err
	}
	return d, nil
}

// Open attempts to connect to the database.
func (s *Source) Open(connURL db.ConnectionURL) error {
	s.connURL = connURL
	return s.open()
}

func (s *Source) open() error {
	_db, err := mongo.Connect(context.Background(),
		options.Client().SetTimeout(connTimeout).ApplyURI(s.connURL.String()))
	if err != nil {
		return err
	}

	err = _db.Ping(context.Background(), nil)
	if err != nil {
		return err
	}

	s.client = _db

	//s.database = db.Database("", nil)

	return nil
}
