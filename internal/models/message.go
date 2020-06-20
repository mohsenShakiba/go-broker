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

	m.Id = string(id)

	// read route
	route, err := r.ReadSlice('\n')

	if err != nil {
		return err
	}

	m.Route = string(route)

	// read payload size
	bSize := make([]byte, 8)
	_, err = r.Read(bSize)

	if err != nil {
		return err
	}

	size := binary.BigEndian.Uint64(bSize)

	// read payload
	bPayload := make([]byte, size)
	_, err = io.ReadFull(r, bPayload)

	if err != nil {
		return err
	}

	m.Payload = bPayload

	return nil
}

// Write will write the message to writer
func (m *Message) Write(w io.Writer) error {

	WriteStr(w, "PUB", m.Id, m.Route)

	// write message type
	_, err := w.Write([]byte("PUB"))

	if err != nil {
		return err
	}

	// write msg id
	_, err = w.Write([]byte(m.Id))

	if err != nil {
		return err
	}

	// write route
	_, err = w.Write([]byte(m.Route))

	if err != nil {
		return err
	}

	// write payload size
	bSize := make([]byte, 8)
	binary.BigEndian.PutUint64(bSize, uint64(len(m.Payload)))

	_, err = w.Write(bSize)

	if err != nil {
		return err
	}

	// write route
	_, err = w.Write(m.Payload)

	if err != nil {
		return err
	}

	return nil
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
