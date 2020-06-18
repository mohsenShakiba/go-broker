package models

import "bufio"

type Err struct {
	Id string
}

func (e *Err) Write(r *bufio.Writer) error {

	// write message type
	_, err := r.Write([]byte("ERR"))

	if err != nil {
		return err
	}

	// write msg id
	_, err = r.Write([]byte(e.Id))

	if err != nil {
		return err
	}

	return nil
}
