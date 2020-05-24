package serializer

import "fmt"

type LineSeparatedSerializer struct {
	Bytes [][]byte
}

func NewLineSeparatedSerializer() *LineSeparatedSerializer {
	return &LineSeparatedSerializer{
		Bytes: make([][]byte, 0, 8),
	}
}

func (ls *LineSeparatedSerializer) WriteStr(key string, value string) {
	ls.Bytes = append(ls.Bytes, []byte(key))
	ls.Bytes = append(ls.Bytes, []byte(":"))
	ls.Bytes = append(ls.Bytes, []byte(value))
	ls.Bytes = append(ls.Bytes, []byte("\n"))
}

func (ls *LineSeparatedSerializer) WriteBytes(key string, value []byte) {
	ls.Bytes = append(ls.Bytes, []byte(key))
	ls.Bytes = append(ls.Bytes, []byte(":"))
	ls.Bytes = append(ls.Bytes, value)
	ls.Bytes = append(ls.Bytes, []byte("\n"))
}

func (ls *LineSeparatedSerializer) GetMessagePrefix() string {

	byteCount := 0

	for _, b := range ls.Bytes {
		byteCount += len(b)
	}

	return fmt.Sprintf("%04d\n", byteCount+5)
}

func (ls *LineSeparatedSerializer) GetMessageBytes() []byte {

	byteCount := 0
	msg := make([]byte, 0)

	for _, b := range ls.Bytes {
		byteCount += len(b)
		msg = append(msg, b...)
	}

	prefix := fmt.Sprintf("%04d\n", byteCount+5)

	return append([]byte(prefix), msg...)
}
