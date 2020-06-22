package models

import (
	"bufio"
	"errors"
	"fmt"
)

func Parse(r *bufio.Reader) (interface{}, error) {
	t1, err := r.ReadSlice('\n')

	if len(t1) < 3 {
		return nil, errors.New(fmt.Sprintf("invalid type: %s", string(t1)))
	}

	// trim the /n
	t := t1[:len(t1)-1]

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
	case "FAKE":
		f := &Ping{}
		err := f.FromReader(r)
		return f, err
	}

	return nil, errors.New(fmt.Sprintf("invalid type: %s", string(t1)))
}
