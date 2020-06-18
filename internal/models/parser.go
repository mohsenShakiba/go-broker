package models

import (
	"bufio"
	"errors"
)

func Parse(r *bufio.Reader) (interface{}, error) {
	t, err := r.ReadSlice('\n')

	if err != nil {
		return nil, err
	}

	switch string(t) {
	case "PUB":
		m := &Message{}
		err := m.FromReader(r)
		return m, err
	case "SUB":
		m := &Register{}
		err := m.FromReader(r)
		return m, err
	case "ACK":
		m := &Ack{}
		err := m.FromReader(r)
		return m, err
	case "NACK":
		m := &Nack{}
		err := m.FromReader(r)
		return m, err
	}

	return nil, errors.New("invalid message type")
}
