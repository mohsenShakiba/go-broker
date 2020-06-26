package storage

import (
	"fmt"
	"go-broker/internal/storage/file"
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

type StorageConfig struct {
	Path        string
	Type        string
	MaxFileSize int64
}

const (
	File   = "F"
	Memory = "M"
)

func NewStorage(conf StorageConfig) Storage {
	switch conf.Type {
	case File:
		return file.NewFileStorage(conf.Path, conf.MaxFileSize)
	case Memory:
		return memory.NewMemoryStorage()
	}
	msg := fmt.Sprintf("Invalid storage type %s", conf.Type)
	panic(msg)
}
