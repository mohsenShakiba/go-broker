package socketserver

import (
	"strconv"
	"testing"
)

// test formatter
func TestFormatter(t *testing.T) {

	input1 := "t1"
	input2 := "t2"

	res := format([]byte(input1), []byte(input2))

	length := string(res[:4])

	l, _ := strconv.Atoi(length)

	actualLength := len(input1)*2 + 1 + messageLengthSize

	if l != actualLength {
		t.Fatalf("the length prefix is invalid length: %d, valid length: %d", l, actualLength)
	}

}
