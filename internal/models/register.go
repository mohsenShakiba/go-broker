package models

import (
	"bufio"
	"encoding/binary"
	"io"
	"strings"
)

type Register struct {
	Id     string
	Routes []string
	Dop    int
}

func (reg *Register) FromReader(r *bufio.Reader) error {

	// read the id
	id, err := r.ReadSlice('\n')

	if err != nil {
		return err
	}

	// trim the /n
	id = id[:len(id)-1]

	reg.Id = string(id)

	// read routes
	routes, err := r.ReadSlice('\n')

	if err != nil {
		return err
	}

	// trim the /n
	routes = routes[:len(routes)-1]

	reg.Routes = strings.Split(string(routes), ",")

	// read dop
	bSize := make([]byte, 4+1)
	_, err = r.Read(bSize)

	if err != nil {
		return err
	}

	size := binary.BigEndian.Uint32(bSize[:4])

	reg.Dop = int(size)

	return nil
}

func (reg *Register) Write(w io.Writer) error {
	err := WriteStr(w, "SUB", reg.Id, strings.Join(reg.Routes, ","))

	if err != nil {
		return err
	}

	return WriteUInt32(w, uint32(reg.Dop))

}
