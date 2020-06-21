package models

import (
	"io"
)

type Err struct {
	Id  string
	Err string
}

func (e *Err) Write(w io.Writer) error {
	return WriteStr(w, "ERR", e.Id, e.Err)
}
