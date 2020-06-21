package models

import (
	"encoding/binary"
	"io"
)

var NewLineDelimiter = []byte("\n")

func WriteStr(w io.Writer, s ...string) error {
	for _, str := range s {
		_, err := w.Write([]byte(str))
		if err != nil {
			return err
		}
		_, err = w.Write(NewLineDelimiter)
		if err != nil {
			return err
		}
	}
	return nil
}

func WriteByte(w io.Writer, bs ...[]byte) error {
	for _, b := range bs {

		// write payload size
		bSize := make([]byte, 8)
		binary.BigEndian.PutUint64(bSize, uint64(len(b)))

		_, err := w.Write(bSize)
		if err != nil {
			return err
		}
		_, err = w.Write(b)
		if err != nil {
			return err
		}
		_, err = w.Write(NewLineDelimiter)
		if err != nil {
			return err
		}
	}
	return nil
}

func WriteUInt32(w io.Writer, i ...uint32) error {
	for _, i32 := range i {
		b := make([]byte, 4)
		binary.BigEndian.PutUint32(b, i32)
		_, err := w.Write(b)
		if err != nil {
			return err
		}
		_, err = w.Write(NewLineDelimiter)
		if err != nil {
			return err
		}
	}
	return nil
}
