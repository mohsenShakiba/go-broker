package memory

import (
	"go-broker/internal/storage"
	"sync"
)

type memoryStorage struct {
	// lock
	l sync.RWMutex
	// payload map
	mp map[string][]byte
}

func NewMemoryStorage() storage.Storage {
	return &memoryStorage{
		mp: make(map[string][]byte),
	}
}

func (ms *memoryStorage) Init() error {
	return nil
}

func (ms *memoryStorage) Keys() ([]string, error) {
	ms.l.Lock()
	defer ms.l.Unlock()
	keys := make([]string, 0, len(ms.mp))
	for id := range ms.mp {
		keys = append(keys, id)
	}
	return keys, nil
}

func (ms *memoryStorage) Write(key string, payload []byte) error {
	ms.l.Lock()
	defer ms.l.Unlock()
	ms.mp[key] = payload
	return nil
}

func (ms *memoryStorage) Read(key string) ([]byte, error) {
	ms.l.RLock()
	defer ms.l.RUnlock()

	b, ok := ms.mp[key]

	if !ok {
		return nil, storage.NotFoundError
	}

	return b, nil
}

func (ms *memoryStorage) Delete(key string) error {
	ms.l.Lock()
	defer ms.l.Unlock()
	delete(ms.mp, key)
	return nil
}

func (ms *memoryStorage) Close() {
}
