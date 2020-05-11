package socketserver

import (
	"fmt"
)

const (
	messageLengthSize   = 4
	messageLengthFormat = "%04d\n"
)

func formatStr(part string) []byte {
	return format([]byte(part))
}

func format(parts ...[]byte) []byte {
	msg := make([]byte, 0, 1024)

	for i, p := range parts {
		msg = append(msg, p...)
		if i < len(parts)-1 {
			msg = append(msg, []byte("\n")...)
		}
	}

	return msg
}

func formatWithLengthPrefix(parts ...[]byte) []byte {

	msg := format(parts...)

	prefix := []byte(fmt.Sprintf(messageLengthFormat, len(msg)+1+messageLengthSize))

	msg = append(prefix, msg...)

	return msg
}
