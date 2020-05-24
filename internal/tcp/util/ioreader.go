package util

import (
	log "github.com/sirupsen/logrus"
	"io"
	"strconv"
)

func ReadLength(b []byte) (int, error) {
	msgPartCountStr := string(b)
	return strconv.Atoi(msgPartCountStr)
}

// Read will read from io based on the message length prefix
func Read(i io.Reader) ([]byte, bool) {

	msgLength, ok := readMessageLength(i)

	if !ok {
		log.Errorf("the first 4 bytes of message isn't int, discarding")
		return nil, false
	}

	b := make([]byte, msgLength)
	_, err := i.Read(b)

	if err != nil {
		log.Errorf("could not read from socket, discarding")
		return nil, false
	}

	return b, true
}

// readMessageLength will determine the length of a message
func readMessageLength(i io.Reader) (int, bool) {

	const messageLength = 5

	// read four bytes
	b := make([]byte, messageLength)
	bytesRead, err := i.Read(b)

	if err != nil {
		return 0, false
	}

	if bytesRead < messageLength {
		return 0, false
	}

	sizeStr := string(b[:messageLength-1])

	length, err := strconv.Atoi(sizeStr)

	if err != nil {
		return 0, false
	}

	return length - messageLength, true
}
