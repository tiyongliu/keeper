package persist

import (
	"errors"
	"keeper/app/db"
	"keeper/app/db/drivers"
	"keeper/app/pkg/logger"
	"keeper/app/pkg/serializer"
	"keeper/app/utility"
	"sync"
)

var lookupIdOnce sync.Once
var lookupIdSession *StorageSession

type repositoryId string
type databaseId string

type StorageSession struct {
	source map[repositoryId]map[databaseId]db.Session
	mu     sync.RWMutex
}

func GetStorageSession() *StorageSession {
	lookupIdOnce.Do(func() {
		lookupIdSession = &StorageSession{
			source: make(map[repositoryId]map[databaseId]db.Session),
		}
	})

	return lookupIdSession
}

func (s *StorageSession) Scanner(conid string, connection map[string]interface{}) (db.Session, error) {
	if conid == "" {
		return nil, db.ErrNilRecord
	}

	var database string
	if connection["database"] != nil {
		database = connection["database"].(string)
	}
	session, err := s.GetItem(conid, database)
	if err != nil || session.Ping() != nil {
		if connection == nil {
			return nil, db.ErrNotConnected
		}

		session, err = drivers.NewCompatDriver().Open(connection)
		if err != nil {
			return nil, err
		}

		if err = s.SetItem(conid, database, session); err != nil {
			return nil, err
		}
	}

	if session == nil {
		return nil, db.ErrNilRecord
	}

	return session, nil
}

func (s *StorageSession) SetItem(conid, database string, driver db.Session) error {
	if driver == nil {
		return errors.New(serializer.ErrNil)
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.source[repositoryId(conid)] == nil {
		s.source[repositoryId(conid)] = make(map[databaseId]db.Session)
	}

	s.source[repositoryId(conid)][databaseId(database)] = driver

	return nil
}

func (s *StorageSession) GetDatabaseMap(conid string) (map[databaseId]db.Session, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	sessionMap, ok := s.source[repositoryId(conid)]
	if !ok {
		return nil, errors.New(serializer.ErrNil)
	}

	return sessionMap, nil
}

func (s *StorageSession) GetItem(conid, database string) (driver db.Session, err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	sessionMap, ok := s.source[repositoryId(conid)]
	if !ok {
		return nil, errors.New(serializer.ErrNil)
	}

	for k, v := range sessionMap {
		if k == databaseId(database) {
			return v, nil
		}
	}

	return nil, errors.New("invalid database")
}

func (s *StorageSession) RemoveItem(conid string) (err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	utility.WithRecover(func() {
		if err = s.closeDriver(conid); err == nil {
			delete(s.source, repositoryId(conid))
		}
	}, func(e error) {
		logger.Errorf("delete driver id failed %v", err)
		err = e
	})

	return err
}

func (s *StorageSession) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, sessions := range s.source {
		for _, session := range sessions {
			session.Close()
		}
	}
	s.source = make(map[repositoryId]map[databaseId]db.Session)
}

func (s *StorageSession) closeDriver(conid string) error {
	sessions := s.source[repositoryId(conid)]
	if sessions != nil {
		for _, session := range sessions {
			session.Close()
		}
		s.source[repositoryId(conid)] = nil
	}
	return nil
}
