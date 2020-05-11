package packer

import (
	"errors"
	"strconv"
)

// asciiPacker will use the ascii representation of number as the prefix
type asciiPacker struct {
}

func (a *asciiPacker) Pack(p []byte) []byte  {

	if p == nil {
		return nil
	}

	messageLength := len(p)

	asciiRepresentation := string(messageLength)

	return append([]byte(asciiRepresentation), p...)
}

func (a *asciiPacker) Unpack(p []byte) ([]byte, error)  {

	if p == nil {
		return nil, errors.New("payload in null")
	}

	if len(p) < 4 {
		return nil, errors.New("payload length does not allow unpacking")
	}

	prefix := string(p[:4])

	length, err := strconv.Atoi(prefix)

	if err != nil {
		return nil, err
	}

	if

	return p[4:], nil

}