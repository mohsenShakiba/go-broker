package models

import (
	"bufio"
	"io"
)

type Ping struct {
	Id string
}

func (p *Ping) FromReader(r *bufio.Reader) error {

	// read the id
	id, err := r.ReadSlice('\n')

	if err != nil {
		return err
	}

	// trim the /n
	id = id[:len(id)-1]

	p.Id = string(id)

	return nil
}

func (p *Ping) Write(w io.Writer) error {
	return WriteStr(w, "PING", p.Id)
}
