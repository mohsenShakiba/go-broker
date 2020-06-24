package storage

import (
	"encoding/binary"
	"go-broker/internal/storage/file"
	"io"
	"os"
	"sync"
)

type fileStorage struct {
	// file path
	path string
	// file handler
	fh *os.File
	// lock for making sure the offset, map and file handler is accessed by only one goroutine
	l sync.RWMutex
	// map containing the file entries in the memory
	m map[string]*file.entry
	// offset in the file for writing new entries
	fhOffset int64
	// offset of actual rows that are not deleted
	fhTrueOffset int64
}

func NewFileStore(path string) Storage {
	return &fileStorage{
		path: path,
		m:    make(map[string]*file.entry),
	}
}

// Init will create a file if it doesn't exists
// and parse the file if it exists
func (fs *fileStorage) Init() error {

	fs.l.Lock()
	defer fs.l.Unlock()

	// open file for read
	fh, err := os.OpenFile(fs.path, os.O_CREATE|os.O_RDWR, 0666)

	if err != nil {
		return err
	}

	fs.fh = fh

	entryByteArr := make([]byte, file.entryHeaderLength)

	for {
		_, err := fh.ReadAt(entryByteArr, fs.fhOffset)

		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		entry := file.fromBinary(entryByteArr)

		entry.offset = fs.fhOffset
		fs.fhOffset += file.entryHeaderLength + entry.length

		if entry.deleted == 0 {
			fs.m[entry.id] = entry
			fs.fhTrueOffset += file.entryHeaderLength + entry.length
		}
	}

	return nil
}

func (fs *fileStorage) Keys() ([]string, error) {
	keys := make([]string, 0, len(fs.m))
	for id := range fs.m {
		keys = append(keys, id)
	}
	return keys, nil
}

func (fs *fileStorage) Read(key string) ([]byte, error) {
	fs.l.RLock()
	defer fs.l.RUnlock()

	// find the entry
	e, ok := fs.m[key]

	if !ok {
		return nil, NotFoundError
	}

	// get the payload
	b := make([]byte, e.length)
	_, err := fs.fh.ReadAt(b, e.offset+20)

	if err != nil {
		return nil, err
	}

	return b, nil
}

func (fs *fileStorage) Write(key string, payload []byte) error {
	fs.l.Lock()
	defer fs.l.Unlock()

	e := &file.entry{
		deleted: 0,
		id:      key,
		length:  int64(len(payload)),
		offset:  fs.fhOffset,
	}

	fs.m[key] = e

	b := file.toBinary(e)

	_, err := fs.fh.WriteAt(b, fs.fhOffset)

	if err != nil {
		return err
	}

	_, err = fs.fh.WriteAt(payload, fs.fhOffset+file.entryHeaderLength)

	if err != nil {
		return err
	}

	fs.fhOffset += int64(len(b) + len(payload))

	return nil
}

func (fs *fileStorage) Delete(key string) error {
	fs.l.Lock()
	defer fs.l.Unlock()

	// find the entry
	e, ok := fs.m[key]

	if !ok {
		return nil
	}

	b := make([]byte, 2)
	binary.LittleEndian.PutUint16(b, 1)

	delete(fs.m, key)

	_, err := fs.fh.WriteAt(b, e.offset)

	fs.fhTrueOffset -= e.length + file.entryHeaderLength

	return err
}

func (fs *fileStorage) Close() {
	fs.l.Lock()
	defer fs.l.Unlock()
	fs.fh.Close()
}

// this method will check if fhOffset is twice as fhTrueOffset
// in which case a new file is created and all the existing data will move to the new file
func (fs *fileStorage) checkForDefragmentation() {

}
