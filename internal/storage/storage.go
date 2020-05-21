package storage

import (
	"fmt"
	"os"
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

	indexPath := fmt.Sprintf("%s/%s", basePath, indexName)
	bodyPath := fmt.Sprintf("%s/%s", basePath, mainName)

	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		os.Create(indexPath)
	}

	if _, err := os.Stat(bodyPath); os.IsNotExist(err) {
		os.Create(bodyPath)
	}

	storage := Storage{
		Path: basePath,
		index: index{
			indices:       make([]*indexRow, 0),
			indexFilePath: indexPath,
			bodyFilePath:  bodyPath,
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

	return s.index.writeAll()
}

func (s *Storage) read(msgId string) *Message {

	// find the index
	var index *indexRow

	for _, i := range s.index.indices {
		if i.msgId == msgId {
			index = i
		}
	}

	if index == nil {
		return nil
	}

	// read body

	body := s.index.readBody(index)

	if body == nil {
		return nil
	}

	return &Message{
		MsgId:   msgId,
		Routes:  body.routes,
		Payload: body.payload,
	}
}
