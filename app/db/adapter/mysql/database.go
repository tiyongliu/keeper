package mysql

import (
	"context"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"keeper/app/db"
	"log"
	"os"
	"sync"
	"time"
)

// Adapter is the public name of the adapter.
const Adapter = `mysql`

type Source struct {
	db.Settings
	ctx            context.Context
	connURL        db.ConnectionURL
	lookupNameOnce sync.Once
	name           string
	mu             sync.Mutex // guards ctx, txOptions
	sqlDBMu        sync.Mutex // guards sess, baseTx
	sqlDB          *gorm.DB
	sessID         uint64
}

func (mysqlAdapter) Open(dsn db.ConnectionURL) (db.Session, error) {
	return Open(dsn)
}

type mysqlAdapter struct {
}

func init() {
	db.RegisterAdapter(Adapter, db.Adapter(&mysqlAdapter{}))
}

func Open(dsn db.ConnectionURL) (db.Session, error) {
	d := &Source{Settings: db.NewSettings(), ctx: context.Background()}
	if err := d.Open(dsn); err != nil {
		return nil, err
	}
	return d, nil
}

func (s *Source) Open(connURL db.ConnectionURL) error {
	s.connURL = connURL
	return s.open()
}

func (s *Source) open() error {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,       // Disable color
		},
	)

	_db, err := gorm.Open(mysql.New(mysql.Config{
		DriverName: Adapter,
		DSN:        s.connURL.String(),
	}), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return err
	}
	s.sqlDB = _db
	return nil
}
