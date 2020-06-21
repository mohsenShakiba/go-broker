package models

import (
	"bufio"
	"io"
)

type Ack struct {
	Id string
}

func (a *Ack) FromReader(r *bufio.Reader) error {

	// read the id
	id, err := r.ReadSlice('\n')

	if err != nil {
		return err
	}

	// trim the /n
	id = id[:len(id)-1]

	a.Id = string(id)

	return nil
}

func (a *Ack) Write(w io.Writer) error {
	return WriteStr(w, "ACK", a.Id)
}
