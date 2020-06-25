package file

import (
	"io/ioutil"
	"strconv"
	"sync"
)

type mapper struct {
	lock            sync.RWMutex
	maxFileSize     int64
	basePath        string
	dataFiles       []*dataFile
	currentDataFile *dataFile
}

func newMapper(basePath string, maxFileSize int64) (*mapper, error) {
	mapper := &mapper{
		maxFileSize:     maxFileSize,
		basePath:        basePath,
		dataFiles:       make([]*dataFile, 0),
		currentDataFile: nil,
	}

	// read folder files
	fi, err := ioutil.ReadDir(basePath)

	if err != nil {
		return nil, err
	}

	// for each file in the folder
	for _, f := range fi {
		name := f.Name()

		index, err := strconv.Atoi(name)

		if err != nil {
			return nil, err
		}

		df, err := openDataFile(basePath, index)

		if err != nil {
			return nil, err
		}

		mapper.dataFiles = append(mapper.dataFiles, df)

		if mapper.currentDataFile != nil && mapper.currentDataFile.index < index {
			mapper.currentDataFile = df
		}
	}

	return mapper, nil
}

func (m *mapper) dataFileForKey(key string) *dataFile {
	m.lock.RLock()
	defer m.lock.RUnlock()
	for _, df := range m.dataFiles {
		if df.containsKey(key) {
			return df
		}
	}
	return nil
}

func (m *mapper) activeDataFile() (*dataFile, error) {

	// if current datafile is still active, return it
	if m.currentDataFile != nil && m.currentDataFile.offset < m.maxFileSize {
		return m.currentDataFile, nil
	}

	m.lock.Lock()
	defer m.lock.Unlock()

	currentIndex := 1

	if m.currentDataFile != nil {
		currentIndex = m.currentDataFile.index + 1
	}

	// create new datafile
	datafile, err := newDataFile(m.basePath, currentIndex)

	if err != nil {
		return nil, err
	}

	m.dataFiles = append(m.dataFiles, datafile)

	m.currentDataFile = datafile
}
