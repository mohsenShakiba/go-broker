package storage

import (
	"fmt"
	"os"
	"sync"
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

}

type storage struct {
	config   StorageConfig
	entryMap map[uint64]entry
}
