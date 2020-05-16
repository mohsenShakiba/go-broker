package socketserver

import (
	"bytes"
	"errors"
	log "github.com/sirupsen/logrus"
	"go-broker/internal/socketserver/serializer"
	"io"
	"strconv"
)

type MessageContext struct {
	PayloadMap map[string][]byte
	socketConn io.ReadWriteCloser
	Serializer serializer.LineSeparatedSerializer
}

func convertToMessage(b []byte) *MessageContext {

	newLineB := []byte("\n")
	colonB := []byte(":")

	tcpMsg := &MessageContext{}

	// split by new line
	partsByNewLine := bytes.Split(b, newLineB)

	for _, part := range partsByNewLine {
		partsByColon := bytes.Split(part, colonB)

		if len(partsByColon) != 2 {
			log.Warnf("bad payload data, discarding, message: %s", string(part))
		}

		tcpMsg.PayloadMap[string(part[0])] = partsByColon[1]
	}

	return tcpMsg
}

func (m *MessageContext) readByteArr(key string) ([]byte, bool) {
	value := m.PayloadMap[key]

	if value == nil {
		return nil, false
	}

	return value, true
}

func (m *MessageContext) readStr(key string) (string, bool) {
	value := m.PayloadMap[key]

	if value == nil {
		return "", false
	}

	return string(value), true
}

func (m *MessageContext) readInt(key string) (int, bool) {
	value := m.PayloadMap[key]

	if value == nil {
		return 0, false
	}

	num, err := strconv.Atoi(string(value))

	if err != nil {
		return 0, false
	}

	return num, true
}

func (m *MessageContext) GetMessageId() (string, bool) {
	return m.readStr("msgId")
}

func (m *MessageContext) GetMessageType() (string, bool) {
	return m.readStr("type")
}

func (m *MessageContext) SendAck() error {

	msgId, ok := m.GetMessageId()

	if !ok {
		return errors.New("the message id isn't provided")
	}

	m.Serializer.WriteStr("type", "ACK")
	m.Serializer.WriteStr("msgId", msgId)

	_, err := m.socketConn.Write([]byte(m.Serializer.GetMessagePrefix()))

	if err != nil {
		return errors.New("failed tow write message prefix to socket connection")
	}

	_, err = m.socketConn.Write(m.Serializer.Bytes)
}

func (m *MessageContext) Close() {
	_ = m.socketConn.Close()
}
