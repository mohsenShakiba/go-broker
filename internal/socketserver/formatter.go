package socketserver

import (
	"fmt"
)

const (
	messageLengthSize   = 4
	messageLengthFormat = "%04d"
)

// this message will convert the message body into []byte
func messageLengthPrefixFormatter(len int) []byte {
	length := fmt.Sprintf(messageLengthFormat, len+messageLengthSize)
	return []byte(length)
}

func format(parts ...[]byte) []byte {
	msg := make([]byte, 0, 1000)

	for _, p := range parts {
		msg = append(msg, p...)
		msg = append(msg, []byte("\n")...)
	}

	msg = append(messageLengthPrefixFormatter(len(msg)), msg...)

	return msg
}
