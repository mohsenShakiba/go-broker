package storage

import (
	"time"
)

type Storage interface {
	Init() error
	Write(id int64, payload []byte) error
	Read(id int64) ([]byte, error)
	Delete(id int64) error
}

type StorageConfig struct {
	Path           string
	SyncPeriod     time.Duration
	FileMaxSize    int64
	FIleNamePrefix string
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
}

func (s *storage) Init() error {

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

func (s *storage) Read(id int64) ([]byte, error) {
	e := s.entryMap[id]

	if e

	return s.handler.readPayload()
}