package file

import (
	"encoding/binary"
	"errors"
	"io"
	"os"
	"sync"
)

var (
	NotFoundErr = errors.New("NotFound")
)

// dataFile is an abstraction over a system file
type dataFile struct {

	// path to backing file
	path string

	// file handler obviously
	f *os.File

	// the current offset in the file
	offset int64

	// the number of entries that have not been expired yet
	remainingActiveEntries int

	// lock to prevent multiple access to file handler
	lock sync.RWMutex
}

// newDataFile will create a new dataFile
func newDataFile(path string) (*dataFile, error) {

	// open file
	fh, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0666)

	if err != nil {
		return nil, err
	}

	// read the size of file
	fileStat, err := fh.Stat()

	if err != nil {
		return nil, err
	}

	df := &dataFile{
		path:   path,
		f:      fh,
		offset: fileStat.Size(),
	}

	return df, nil
}

func (df *dataFile) readEntries() ([]*entry, error) {

	// byte array for reading the header
	entries := make([]*entry, 0, 1024)
	entryByteArr := make([]byte, entryHeaderLength)

	// process all the entries
	for {

		// read the header
		_, err := df.f.ReadAt(entryByteArr, df.offset)

		if err == io.EOF {
			break
		}

		if err != nil {
			return entries, err
		}

		// parse from binary
		entry := fromBinary(entryByteArr)

		// adjust parameters
		entry.df = df

		// adjust the offset
		entry.offset = df.offset
		df.offset += entryHeaderLength + entry.length

		// add to the map if not deleted
		if entry.deleted == 0 {
			entries = append(entries, entry)
			df.remainingActiveEntries += 1
		}
	}

	return entries, nil
}

// append will add a new entry to the current file
func (df *dataFile) append(e *entry, value []byte) error {
	df.lock.Lock()
	defer df.lock.Unlock()

	// write the header
	_, err := df.f.WriteAt(toBinary(e), df.offset)

	if err != nil {
		return err
	}

	// write the payload
	_, err = df.f.WriteAt(value, df.offset+entryHeaderLength)

	// adjust the offset
	df.offset += df.offset + entryHeaderLength

	// adjust active entries
	df.remainingActiveEntries += 1

	return nil
}

func (df *dataFile) read(e *entry) ([]byte, error) {
	df.lock.RLock()
	defer df.lock.RUnlock()

	b := make([]byte, e.length)
	_, err := df.f.ReadAt(b, e.offset+entryHeaderLength)

	if err != nil {
		return nil, err
	}

	return b, nil
}

func (df *dataFile) delete(e *entry) error {
	df.lock.Lock()
	defer df.lock.Lock()

	b := make([]byte, 2)
	binary.LittleEndian.PutUint16(b, 1)

	// change 0 to 1 indicating that the entry has been deleted
	_, err := df.f.WriteAt(b, e.offset)

	// adjust active entries
	df.remainingActiveEntries += 1

	return err
}

func (df *dataFile) close() {
	_ = df.f.Close()
}
