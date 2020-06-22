package models

import "bufio"

type Ping struct {
	Id string
}

func (f *Ping) FromReader(r *bufio.Reader) error {

	// read the id
	id, err := r.ReadSlice('\n')

	if err != nil {
		return err
	}

	// trim the /n
	id = id[:len(id)-1]

	f.Id = string(id)

	return nil
}
