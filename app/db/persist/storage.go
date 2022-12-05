package persist

import (
	"errors"
	"keeper/app/db"
	"keeper/app/db/drivers"
	"keeper/app/pkg/logger"
	"keeper/app/utility"
	"sync"
)

var lookupIdOnce sync.Once
var lookupIdSession *StorageSession

type StorageSession struct {
	source map[string]db.Session
	mu     sync.RWMutex
}

func GetStorageSession() *StorageSession {
	lookupIdOnce.Do(func() {
		lookupIdSession = &StorageSession{
			source: make(map[string]db.Session),
		}
	})

	return lookupIdSession
}

func (s *StorageSession) Scanner(conid string, connection map[string]interface{}) (db.Session, error) {
	if conid == "" {
		return nil, db.ErrNilRecord
	}

	session, err := s.GetItem(conid)
	if err != nil {
		if connection == nil {
			return nil, db.ErrNotConnected
		}

		session, err = drivers.NewCompatDriver().Open(connection)
		if err != nil {
			return nil, err
		}

		if err = s.SetItem(conid, session); err != nil {
			return nil, err
		}
	}

	if session == nil {
		return nil, db.ErrNilRecord
	}

	return session, nil
}

func (s *StorageSession) SetItem(conid string, driver db.Session) error {
	if driver == nil {
		return errors.New("invalid memory address or nil pointer dereference")
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	s.source[conid] = driver
	return nil
}

func (s *StorageSession) GetItem(conid string) (driver db.Session, err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	session, ok := s.source[conid]
	if !ok {
		return nil, errors.New("invalid memory address or nil pointer dereference")
	}

	return session, nil
}

func (s *StorageSession) RemoveItem(conid string) (err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	utility.WithRecover(func() {
		if err = s.closeDriver(conid); err == nil {
			delete(s.source, conid)
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
	for _, session := range s.source {
		session.Close()
	}
	s.source = make(map[string]db.Session)
}

func (s *StorageSession) closeDriver(conid string) error {
	session := s.source[conid]
	if session != nil {
		return session.Close()
	}
	return nil
}
