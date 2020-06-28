package file

import (
	"fmt"
	"io/ioutil"
	"path"
	"sync"
	"time"
)

type fileStorage struct {
	msgMap         map[string]*entry
	dataFiles      []*dataFile
	activeDataFile *dataFile
	basePath       string
	maxFileSize    int64
	lock           sync.RWMutex
}

func NewFileStorage(basePath string, maxSize int64) *fileStorage {
	return &fileStorage{
		msgMap:         make(map[string]*entry),
		dataFiles:      make([]*dataFile, 0),
		activeDataFile: nil,
		basePath:       basePath,
		maxFileSize:    maxSize,
		lock:           sync.RWMutex{},
	}
}

func (fs *fileStorage) Init() error {
	// read folder files
	fi, err := ioutil.ReadDir(fs.basePath)

	if err != nil {
		return err
	}

	// for each file in the folder
	for _, f := range fi {
		name := f.Name()

		dfPath := path.Join(fs.basePath, name)

		// append data file
		df, err := newDataFile(dfPath)
		fs.dataFiles = append(fs.dataFiles, df)
		fs.activeDataFile = df

		if err != nil {
			return err
		}

		entries, err := df.readEntries()

		if err != nil {
			return err
		}

		for _, e := range entries {
			fs.msgMap[e.id] = e
		}

	}

	return nil
}

func (fs *fileStorage) Keys() ([]string, error) {
	fs.lock.RLock()
	defer fs.lock.RUnlock()

	keys := make([]string, 0, len(fs.msgMap))
	for id := range fs.msgMap {
		keys = append(keys, id)
	}
	return keys, nil
}

func (fs *fileStorage) Read(key string) ([]byte, error) {
	fs.lock.RLock()
	defer fs.lock.RUnlock()

	e, ok := fs.msgMap[key]

	if !ok {
		return nil, NotFoundErr
	}

	return e.df.read(e)
}

func (fs *fileStorage) Write(key string, payload []byte) error {
	// get current dataFile
	dataFile, err := fs.getActiveDataFile()

	if err != nil {
		return err
	}

	e := &entry{
		deleted: 0,
		id:      key,
		length:  int64(len(payload)),
		df:      dataFile,
		offset:  dataFile.offset,
	}

	fs.lock.Lock()
	defer fs.lock.Unlock()

	fs.msgMap[key] = e

	return dataFile.append(e, payload)
}

func (fs *fileStorage) Delete(key string) error {
	fs.lock.Lock()
	defer fs.lock.Unlock()

	// find the entry
	e, ok := fs.msgMap[key]

	if !ok {
		return nil
	}

	delete(fs.msgMap, key)

	return e.df.delete(e)
}

func (fs *fileStorage) Close() {
	fs.lock.Lock()
	defer fs.lock.Unlock()

	for _, df := range fs.dataFiles {
		df.close()
	}

}

func (fs *fileStorage) getActiveDataFile() (*dataFile, error) {
	if fs.activeDataFile != nil && fs.activeDataFile.offset <= fs.maxFileSize {
		return fs.activeDataFile, nil
	}

	fs.lock.Lock()
	defer fs.lock.Unlock()

	dfPath := path.Join(fs.basePath, fmt.Sprintf("%d", time.Now().Unix()))

	df, err := newDataFile(dfPath)

	if err != nil {
		return nil, err
	}

	fs.activeDataFile = df
	fs.dataFiles = append(fs.dataFiles, df)

	return fs.activeDataFile, nil
}
