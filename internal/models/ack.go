package models

import (
	"bufio"
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

	a.Id = string(id)

	return nil
}
