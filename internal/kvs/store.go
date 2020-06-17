package kvs

import (
	"io"
	"os"
	"sync"
)

// Store represent a file which contains key values for a specific channel
type Store interface {
	Init(path string) error
	Keys() ([]string, error)
	Write(k string, v []byte) error
	Read(k string) ([]byte, error)
	Delete(k string) error
	Close()
}

type store struct {
	fh *os.File
	// lock for making sure the offset, map and file handler is accessed by only one goroutine
	l  sync.RWMutex
	// map containing the file entries in the memory
	m map[int64]*entry
	// offset in the file for writing new entries
	fhOffset int64
	// offset of actual rows that are not deleted
	fhTrueOffset int64
}

// Init will create a file if it doesn't exists
// and parse the file if it exists
func (s *store) Init(path string) error {

	fh, err := os.OpenFile(path, os.O_CREATE | os.O_RDWR, 0666)

	if err != nil {
		return err
	}

	var len int64 = 20

	entryByteArr := make([]byte, 20)

	for {
		_, err := fh.ReadAt(entryByteArr, s.fhOffset)

		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		entry := fromBinary(entryByteArr)

		entry.offset = s.fhOffset
		s.fhOffset += len + entry.length

		if entry.deleted == 0 {
			s.m[entry.id] = entry
			s.fhTrueOffset += len + entry.length
		}
	}

	return nil
}

func (s *store) Keys() ([]string, error) {
	keys := make([]string, 0, len(s.m))
	for k
	:= range s.m {
		keys = append(keys, k)
	}
	return keys, nil
}


