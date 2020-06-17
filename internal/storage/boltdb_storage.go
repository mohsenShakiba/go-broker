package storage

import (
	"github.com/boltdb/bolt"
	"sync"
)

type boltDbStorage struct {
	rwMutex sync.RWMutex
	db      *bolt.DB
	path    string
}

func NewBoltDbStorage(path string) Storage {
	return &boltDbStorage{
		path: path,
	}
}

func (b *boltDbStorage) Init() error {
	db, err := bolt.Open(b.path, 0600, nil)

	if err != nil {
		return err
	}

	b.db = db

	return nil
}

func (b *boltDbStorage) Keys() ([]string, error) {
	keys := make([]string, 0)
	err := b.db.View(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("keys"))
		if err != nil {
			return err
		}
		cursor := b.Cursor()
		for k, _ := cursor.First(); k != nil; k, _ = cursor.Next() {
			keys = append(keys, string(k))
		}

		return nil
	})
	return keys, err
}

func (b *boltDbStorage) Write(key string, payload []byte) error {
	err := b.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("keys"))
		if err != nil {
			return err
		}
		return b.Put([]byte(key), payload)
	})
	return err
}

func (b *boltDbStorage) Read(key string) ([]byte, error) {
	var payload []byte
	err := b.db.View(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("keys"))
		if err != nil {
			return err
		}
		payload = b.Get([]byte(key))
		return nil
	})
	return payload, err
}

func (b *boltDbStorage) Delete(key string) error {
	err := b.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("keys"))
		if err != nil {
			return err
		}
		return b.Delete([]byte(key))
	})
	return err
}

func (b *boltDbStorage) Close() {
	_ = b.db.Close()
}
