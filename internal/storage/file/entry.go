package file

import (
	"bytes"
	"encoding/binary"
)

const (
	entryHeaderLength = 32 + 10
)

type entry struct {
	deleted uint16
	id      string
	offset  int64
	length  int64
}

func toBinary(e *entry) []byte {
	b := make([]byte, entryHeaderLength)
	binary.LittleEndian.PutUint16(b[:2], e.deleted)
	binary.LittleEndian.PutUint16(b[2:10], uint16(e.length))
	copy(b[10:], e.id)
	return b
}

func fromBinary(b []byte) *entry {
	return &entry{
		deleted: binary.LittleEndian.Uint16(b[:2]),
		length:  int64(binary.LittleEndian.Uint64(b[2:10])),
		id:      string(bytes.Trim(b[10:], "\x00")),
	}
}
