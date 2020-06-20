package models

import "io"

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
	for _, str := range bs {
		_, err := w.Write(str)
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
