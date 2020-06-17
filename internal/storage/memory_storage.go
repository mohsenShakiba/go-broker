package storage

import (
	"sync"
)

type memoryStorage struct {
	// lock
	l sync.RWMutex
	// entry map
	m map[string]*entry
	// payload map
	mp map[string][]byte
}

func NewMemoryStore() Storage {
	return &memoryStorage{
		m:  make(map[string]*entry),
		mp: make(map[string][]byte),
	}
}

func (ms *memoryStorage) Init() error {
	return nil
}

func (ms *memoryStorage) Keys() ([]string, error) {
	ms.l.Lock()
	defer ms.l.Unlock()
	keys := make([]string, 0, len(ms.m))
	for id := range ms.m {
		keys = append(keys, id)
	}
	return keys, nil
}

func (ms *memoryStorage) Write(key string, payload []byte) error {
	ms.l.Lock()
	defer ms.l.Unlock()
	e := &entry{
		deleted: 0,
		id:      key,
		length:  int64(len(payload)),
	}

	ms.m[key] = e
	ms.mp[key] = payload

	return nil
}

func (ms *memoryStorage) Read(key string) ([]byte, error) {
	ms.l.RLock()
	defer ms.l.RUnlock()

	b, ok := ms.mp[key]

	if !ok {
		return nil, NotFoundError
	}

	return b, nil
}

func (ms *memoryStorage) Delete(key string) error {
	ms.l.Lock()
	defer ms.l.Unlock()

	delete(ms.m, key)
	delete(ms.mp, key)

	return nil
}

func (ms *memoryStorage) Close() {
}
