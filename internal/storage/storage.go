package storage

import (
	"errors"
	"sync"
	"time"
)

type Storage interface {
	Init() error
	Write(id int64, payload []byte) error
	Read(id int64) ([]byte, error)
	Delete(id int64) error
	Dispose()
}

type StorageConfig struct {
	Path           string
	SyncPeriod     time.Duration
	FileMaxSize    int64
	FileNamePrefix string
}

func New(cng StorageConfig) Storage {
	return &storage{
		config:   cng,
		entryMap: make(map[int64]*entry),
	}
}

type storage struct {
	config   StorageConfig
	handler  *fileHandler
	entryMap map[int64]*entry
	lock     sync.Mutex
}

func (s *storage) Init() error {

	s.lock.Lock()
	defer s.lock.Unlock()

	s.handler = newHandler(s.config)

	entries, err := s.handler.readAllEntries()

	if err != nil {
		return err
	}

	for _, entry := range entries {
		s.entryMap[entry.id] = entry
	}

	return nil
}

func (s *storage) Write(id int64, payload []byte) error {
	e, err := s.handler.write(id, payload)
	s.lock.Lock()
	s.entryMap[id] = e
	s.lock.Unlock()
	return err
}

func (s *storage) Read(id int64) ([]byte, error) {
	s.lock.Lock()
	e := s.entryMap[id]
	s.lock.Unlock()

	if e == nil {
		return nil, errors.New("entry not found")
	}

	return s.handler.readPayload(e)
}

func (s *storage) Delete(id int64) error {
	s.lock.Lock()
	e := s.entryMap[id]
	s.lock.Unlock()

	if e == nil {
		return errors.New("entry not found")
	}

	s.lock.Lock()
	delete(s.entryMap, id)
	s.lock.Unlock()

	return s.handler.delete(e)
}

func (s *storage) Dispose() {
	s.handler.dispose()
}
