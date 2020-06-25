package file

import (
	"encoding/binary"
	"io"
	"os"
)

type fileStorage struct {
	m mapper
}

func (fs *fileStorage) Init() error {
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
