package storage

import (
	"fmt"
	"sync"
)

const (
	indexName     = "index.txt"
	indexTempName = "tem_index.txt"
	mainName      = "main.txt"
)

type Message struct {
	MsgId   string
	Routes  []string
	Payload []byte
}

type Storage struct {
	Path   string
	index  index
	msgMap map[string]*indexRow
	mutex  sync.Mutex
}

func Init(basePath string) (*Storage, error) {
	storage := Storage{
		Path: basePath,
		index: index{
			indices:       make([]*indexRow, 0),
			indexFilePath: fmt.Sprintf("%s.%s", basePath, indexName),
			bodyFilePath:  fmt.Sprintf("%s.%s", basePath, mainName),
		},
		msgMap: make(map[string]*indexRow),
		mutex:  sync.Mutex{},
	}

	err := storage.index.readAll()

	if err != nil {
		return nil, err
	}

	return &storage, nil

}

func (s *Storage) Add(m *Message) error {

	offset, length, err := s.index.writeBody(&bodyRow{
		routes:  m.Routes,
		payload: m.Payload,
	})

	if err != nil {
		return err
	}

	indexRow := &indexRow{
		msgId:   m.MsgId,
		deleted: false,
		start:   offset,
		length:  length,
	}

	err = s.index.addIndex(indexRow)

	s.msgMap[m.MsgId] = indexRow

	return err
}

func (s *Storage) Remove(msgId string) error {

	indexRow, ok := s.msgMap[msgId]

	if !ok {
		return nil
	}

	indexRow.deleted = true

	return err
}
