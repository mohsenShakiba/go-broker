package storage

import (
	"io"
	"os"
	"sync"
)

type fileHandler struct {
	path string
	lock sync.RWMutex
}

func (fh *fileHandler) readAllEntries() ([]entry, error) {
	fh.lock.Lock()
	defer fh.lock.Unlock()

	// read the file content
	handler, err := os.Open(fh.path)
	defer handler.Close()

	if err != nil {
		return nil, err
	}

	var off int64 = 0
	b := make([]byte, 0, 18)
	blen := int64(18)

	entries := make([]entry, 0, 124)

	for {
		_, err := handler.ReadAt(b, off)

		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		entry := fromBinary(b)

		entry.offset = off
		off += blen + int64(entry.length)

		entries = append(entries, entry)
	}

	return entries, nil
}

func (fh *fileHandler) readPayload(e entry) ([]byte, error) {
	fh.lock.RLock()
	defer fh.lock.RUnlock()

	// read the file content
	handler, err := os.Open(fh.path)
	defer handler.Close()

	if err != nil {
		return nil, err
	}

	b := make([]byte, e.length)
	_, err = handler.ReadAt(b, e.offset+18)

	if err != nil {
		return nil, err
	}

	return b, nil
}

func (fh *fileHandler) write(id uint64, payload []byte) (entry, error) {

}

func (fh *fileHandler) delete(id uint64) error {

}
