package file

import (
	"bytes"
	"encoding/binary"
)

const (
	entryHeaderLength = 32 + 10
)

// entry contains the information about a key value pair
// storage specific information is stored in the entryStat
type entry struct {

	// the key for entry
	// it's max size is 32 bytes
	id string

	// boolean flag representing if the entry has been deleted
	deleted uint16

	// length of value for this entry
	length int64

	// the offset from start of the file in datafile
	offset int64

	// data file this entry belongs to
	df *dataFile
}

// to binary will create a byte array containing the entry data
func toBinary(e *entry) []byte {
	b := make([]byte, entryHeaderLength)
	binary.LittleEndian.PutUint16(b[:2], e.deleted)
	binary.LittleEndian.PutUint16(b[2:10], uint16(e.length))
	copy(b[10:], e.id)
	return b
}

// from binary will convert a byte array to entry
func fromBinary(b []byte) *entry {
	return &entry{
		deleted: binary.LittleEndian.Uint16(b[:2]),
		length:  int64(binary.LittleEndian.Uint64(b[2:10])),
		id:      string(bytes.Trim(b[10:], "\x00")),
	}
}
