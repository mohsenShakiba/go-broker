package storage

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
)

type index struct {
	indices       []*indexRow
	indexFilePath string
	bodyFilePath  string
}

func (i *index) addIndex(r *indexRow) error {
	i.indices = append(i.indices, r)
	return i.writeAll()
}

func (i *index) readAll() error {
	fileContent, err := ioutil.ReadFile(i.indexFilePath)

	if err != nil {
		return err
	}

	fileContentLineSeparated := bytes.Split(fileContent, []byte("\n"))

	indexRows := make([]*indexRow, 0, len(fileContentLineSeparated))

	for _, indexRowBinary := range fileContentLineSeparated {
		indexRow := deserializeIndex(indexRowBinary)

		if indexRow == nil {
			continue
		}

		indexRows = append(indexRows, indexRow)

	}

	return nil
}

func (i *index) writeAll() error {
	buff := bytes.Buffer{}
	lineSep := []byte("\n")

	for _, indexRow := range i.indices {
		buff.Write(indexRow.serializeIndex())
		buff.Write(lineSep)
	}

	return ioutil.WriteFile(i.indexFilePath, buff.Bytes(), 0644)
}

func (i *index) readBody(r *indexRow) *bodyRow {
	fd, err := os.Open(i.bodyFilePath)
	if err != nil { //error handler
		return nil
	}

	b := make([]byte, r.length)

	_, err = fd.ReadAt(b, int64(r.start))

	if err != nil {
		return nil
	}

	bodyRow := &bodyRow{}

	err = json.Unmarshal(b, bodyRow)

	if err != nil {
		return nil
	}

	return bodyRow
}

func (i *index) writeBody(r *bodyRow) (int, int, error) {
	fd, err := os.OpenFile(i.bodyFilePath, os.O_APPEND|os.O_WRONLY, 0600)

	if err != nil {
		return 0, 0, err
	}

	b, err := json.Marshal(r)

	if err != nil {
		return 0, 0, err
	}

	s, err := fd.Stat()

	if err != nil {
		return 0, 0, err
	}

	_, err = fd.Write(b)

	if err != nil {
		return 0, 0, err
	}

	return int(s.Size()), len(b), nil
}
