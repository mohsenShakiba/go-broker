package socketserver

import (
	"bytes"
	log "github.com/sirupsen/logrus"
	"strconv"
)

type tcpMessage struct {
	payloadMap map[string][]byte
}

func convertToMessage(b []byte) *tcpMessage {

	newLineB := []byte("\n")
	colonB := []byte(":")

	tcpMsg := &tcpMessage{}

	// split by new line
	partsByNewLine := bytes.Split(b, newLineB)

	for _, part := range partsByNewLine {
		partsByColon := bytes.Split(part, colonB)

		if len(partsByColon) != 2 {
			log.Warnf("bad payload data, discarding, message: %s", string(part))
		}

		tcpMsg.payloadMap[string(part[0])] = partsByColon[1]
	}

	return tcpMsg
}

func (m *tcpMessage) readByteArr(key string) ([]byte, bool) {
	value := m.payloadMap[key]

	if value == nil {
		return nil, false
	}

	return value, true
}

func (m *tcpMessage) readStr(key string) (string, bool) {
	value := m.payloadMap[key]

	if value == nil {
		return "", false
	}

	return string(value), true
}

func (m *tcpMessage) readInt(key string) (int, bool) {
	value := m.payloadMap[key]

	if value == nil {
		return 0, false
	}

	num, err := strconv.Atoi(string(value))

	if err != nil {
		return 0, false
	}

	return num, true
}
