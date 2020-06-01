package storage

import "encoding/binary"

type entry struct {
	deleted uint16
	id      uint64
	length  uint64
	offset  int64
}

func toBinary(e entry) []byte {
	totalSize := 8*2 + 2
	b := make([]byte, totalSize)
	binary.LittleEndian.PutUint16(b[:2], e.deleted)
	binary.LittleEndian.PutUint64(b[2:10], e.id)
	binary.LittleEndian.PutUint64(b[10:18], e.length)
	return b
}

func fromBinary(b []byte) entry {
	return entry{
		deleted: binary.LittleEndian.Uint16(b[:2]),
		id:      binary.LittleEndian.Uint64(b[2:10]),
		length:  binary.LittleEndian.Uint64(b[10:18]),
	}
}
