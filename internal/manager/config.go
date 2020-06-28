package manager

import "go-broker/internal/storage"

type Config struct {
	Port int32
	storage.StorageConfig
}
