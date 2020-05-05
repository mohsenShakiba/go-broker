package socketserver

import "fmt"

const (
	messageLengthSize   = 4
	messageLengthFormat = "%04d"
)

// this message will convert the message body into []byte
func messageLengthPrefixFormatter(len int) []byte {
	length := fmt.Sprintf(messageLengthFormat, len+messageLengthSize)
	return []byte(length)
}
