package models

import (
	"bufio"
	"encoding/binary"
	"strings"
)

type Register struct {
	Id     string
	Dop    int
	Routes []string
}

func (m *Register) FromReader(r *bufio.Reader) error {

	// read the id
	id, err := r.ReadSlice('\n')

	if err != nil {
		return err
	}

	m.Id = string(id)

	// read routes
	routes, err := r.ReadSlice('\n')

	if err != nil {
		return err
	}

	m.Routes = strings.Split(string(routes), ",")

	// read dop
	bSize := make([]byte, 2)
	_, err = r.Read(bSize)

	if err != nil {
		return err
	}

	size := binary.BigEndian.Uint64(bSize)

	m.Dop = int(size)

	return nil
}
