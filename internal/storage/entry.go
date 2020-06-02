package storage

import "encoding/binary"

type entry struct {
	deleted int8
	id      int64
	length  int64
	pageNo  int16
	offset  int64
}

func toBinary(e *entry) []byte {
	totalSize := 8*2 + 2
	b := make([]byte, totalSize)
	binary.LittleEndian.PutUint16(b[:2], uint16(e.deleted))
	binary.LittleEndian.PutUint64(b[2:10], uint64(e.id))
	binary.LittleEndian.PutUint64(b[10:18], uint64(e.length))
	binary.LittleEndian.PutUint16(b[18:20], uint16(e.pageNo))
	return b
}

func fromBinary(b []byte) *entry {
	return &entry{
		deleted: int8(binary.LittleEndian.Uint16(b[:2])),
		id:      int64(binary.LittleEndian.Uint64(b[2:10])),
		length:  int64(binary.LittleEndian.Uint64(b[10:18])),
		pageNo:  int16(binary.LittleEndian.Uint16(b[18:20])),
	}
}
