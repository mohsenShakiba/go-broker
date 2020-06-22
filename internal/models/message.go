package models

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"io"
)

// Message message contains the actual payload data and routing data
type Message struct {
	Id      string
	Route   string
	Payload []byte
}

// FromReader will create the message from and io.Reader
func (m *Message) FromReader(r *bufio.Reader) error {
	// read the id
	id, err := r.ReadSlice('\n')

	if err != nil {
		return err
	}

	// trim the /n
	id = id[:len(id)-1]

	m.Id = string(id)

	// read route
	route, err := r.ReadSlice('\n')

	if err != nil {
		return err
	}

	// trim the /n
	route = route[:len(route)-1]

	m.Route = string(route)

	// read payload size
	bSize := make([]byte, 8)
	_, err = io.ReadFull(r, bSize)

	if err != nil {
		return err
	}

	size := binary.BigEndian.Uint64(bSize)

	// read payload
	// +1 for /n
	bPayload := make([]byte, size+1)
	_, err = io.ReadFull(r, bPayload)

	if err != nil {
		return err
	}

	m.Payload = bPayload[:size]

	return nil
}

// Write will write the message to writer
func (m *Message) Write(w io.Writer) error {

	err := WriteStr(w, "PUB", m.Id, m.Route)

	if err != nil {
		return err
	}

	return WriteByte(w, m.Payload)
}

func (m *Message) ToBinary() ([]byte, error) {
	b, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (m *Message) FromBinary(b []byte) error {
	err := json.Unmarshal(b, m)
	if err != nil {
		return err
	}
	return nil
}
