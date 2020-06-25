package file

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"sync"
)

var (
	NotFoundErr = errors.New("NotFound")
)

// dataFile is an abstraction over a system file
type dataFile struct {
	// base path for the containing folder
	path string

	// the index of file in directory
	index int

	// file handler obviously
	f *os.File

	// the current offset in the file
	offset int64

	// the number of entries that have not been expired yet
	remainingActiveEntries int

	// map to know the entry for each key
	entryMap map[string]*entry

	// lock to prevent multiple access to file handler
	lock sync.RWMutex
}

// newDataFile will create a new dataFile
func newDataFile(basePath string, index int) (*dataFile, error) {

	// open file
	fileName := fmt.Sprintf("f%d", index)
	filePath := path.Join(basePath, fileName)
	fh, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, 0666)

	if err != nil {
		return nil, err
	}

	df := &dataFile{
		path:                   basePath,
		index:                  index,
		f:                      fh,
		offset:                 0,
		remainingActiveEntries: 0,
		entryMap:               make(map[string]*entry),
	}

	return df, nil
}

// openDataFile will open an existing dataFile
func openDataFile(basePath string, index int) (*dataFile, error) {

	fileName := fmt.Sprintf("f%d", index)
	filePath := path.Join(basePath, fileName)
	fh, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, 0666)

	if err != nil {
		return nil, err
	}

	// read the size of file
	fileStat, err := fh.Stat()

	if err != nil {
		return nil, err
	}

	// parse the entries

	// byte array for reading the header
	entryByteArr := make([]byte, entryHeaderLength)
	entriesMap := make(map[string]*entry)
	var offset int64
	var activeEntryCount int

	for {

		// read the header
		_, err := fh.ReadAt(entryByteArr, offset)

		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		// parse from binary
		entry := fromBinary(entryByteArr)

		// adjust the offset
		entry.offset = offset
		offset += entryHeaderLength + entry.length

		// add to the map if not deleted
		if entry.deleted == 0 {
			entriesMap[entry.id] = entry
			activeEntryCount += 1
		}
	}

	df := &dataFile{
		path:                   basePath,
		index:                  index,
		f:                      fh,
		offset:                 fileStat.Size(),
		remainingActiveEntries: activeEntryCount,
		entryMap:               entriesMap,
	}

	return df, nil
}

// append will add a new entry to the current file
func (df *dataFile) append(key string, value []byte) error {
	df.lock.Lock()
	defer df.lock.Unlock()

	e := &entry{
		deleted: 0,
		id:      key,
		offset:  df.offset,
		length:  int64(len(value)),
	}

	df.entryMap[key] = e

	// write the header
	_, err := df.f.WriteAt(toBinary(e), df.offset)

	if err != nil {
		return err
	}

	// write the payload
	_, err = df.f.WriteAt(value, df.offset+entryHeaderLength)

	// adjust the offset
	df.offset += df.offset + entryHeaderLength

	return nil
}

func (df *dataFile) read(key string) ([]byte, error) {
	df.lock.RLock()
	defer df.lock.RUnlock()

	// find the entry
	e, ok := df.entryMap[key]

	if !ok {
		return nil, NotFoundErr
	}

	// read at the offset
	b := make([]byte, e.length)
	_, err := df.f.ReadAt(b, e.offset+entryHeaderLength)

	if err != nil {
		return nil, err
	}

	return b, nil
}

func (df *dataFile) delete(key string) error {
	df.lock.Lock()
	defer df.lock.Lock()

	// find the entry
	e, ok := df.entryMap[key]

	if !ok {
		return nil
	}

	b := make([]byte, 2)
	binary.LittleEndian.PutUint16(b, 1)

	delete(df.entryMap, key)

	// change 0 to 1 indicating that the entry has been deleted
	_, err := df.f.WriteAt(b, e.offset)

	return err
}

func (df *dataFile) containsKey(key string) bool {
	df.lock.RLock()
	defer df.lock.RUnlock()

	// find the entry
	_, ok := df.entryMap[key]

	return ok
}
