package models

import "bufio"

type Err struct {
	Id  string
	Err string
}

func (e *Err) Write(r *bufio.Writer) error {
	return WriteStr(r, "ERR", e.Id, e.Err)
}
