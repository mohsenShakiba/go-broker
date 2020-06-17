package models

import "io"

type Ack struct {
	Id string
}

func FromReader(r io.Reader) (*Ack, error) {

}

func ToWriter(r io.Writer) error {

}
