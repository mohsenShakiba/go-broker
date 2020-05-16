package serializer

import "fmt"

type LineSeparatedSerializer struct {
	Bytes []byte
}

func NewLineSeparatedSerializer(size int) *LineSeparatedSerializer {
	return &LineSeparatedSerializer{
		Bytes: make([]byte, 0, size),
	}
}

func (ls *LineSeparatedSerializer) WriteStr(key string, value string) {
	ls.Bytes = append(ls.Bytes, []byte(key)...)
	ls.Bytes = append(ls.Bytes, []byte(":")...)
	ls.Bytes = append(ls.Bytes, []byte(value)...)
	ls.Bytes = append(ls.Bytes, []byte("\n")...)
}

func (ls *LineSeparatedSerializer) WriteBytes(key string, value []byte) {
	ls.Bytes = append(ls.Bytes, []byte(key)...)
	ls.Bytes = append(ls.Bytes, []byte(":")...)
	ls.Bytes = append(ls.Bytes, value...)
	ls.Bytes = append(ls.Bytes, []byte("\n")...)
}

func (ls *LineSeparatedSerializer) GetMessagePrefix() string {
	return fmt.Sprintf("%04d\n", len(ls.Bytes)+7)
}
