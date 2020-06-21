package models

import (
	"bufio"
	"io"
)

type Nack struct {
	Id string
}

func (n *Nack) FromReader(r *bufio.Reader) error {

	// read the id
	id, err := r.ReadSlice('\n')

	if err != nil {
		return err
	}

	// trim the /n
	id = id[:len(id)-1]

	n.Id = string(id)

	return nil
}

func (n *Nack) Write(w io.Writer) error {
	return WriteStr(w, "NACK", n.Id)
}
