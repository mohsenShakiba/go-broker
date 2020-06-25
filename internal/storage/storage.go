package storage

import (
	"errors"
	"fmt"
	"go-broker/internal/storage/memory"
)

// Storage represent a file which contains key values for a specific channel
type Storage interface {
	Init() error
	Keys() ([]string, error)
	Write(k string, v []byte) error
	Read(k string) ([]byte, error)
	Delete(k string) error
	Close()
}

var NotFoundError = errors.New("NOT_FOUND")

const (
	File   = "F"
	Memory = "M"
	BoltDb = "B"
)

func NewStorage(path string, t string) Storage {
	switch t {
	case File:
		return NewFileStore(path)
	case Memory:
		return memory.NewMemoryStore()
	case BoltDb:
		return NewBoltDbStorage(path)
	}
	msg := fmt.Sprintf("Invalid storage type %s", t)
	panic(msg)
}
