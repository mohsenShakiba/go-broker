package models

import (
	"encoding/json"
	"go-broker/internal/tcp/messages"
	"strings"
)

// Message message contains the actual payload data and routing data
type Message struct {
	Id      string
	Routes  []string
	Payload []byte
}

func (m *Message) FromTcpMessage(msg *messages.Message) bool {

	routesStr, ok := msg.ReadStr("routes")

	if !ok {
		return false
	}

	routes := strings.Split(routesStr, ",")

	payload, ok := msg.ReadByteArr("payload")

	if !ok {
		return false
	}

	m.Id = msg.MsgId
	m.Routes = routes
	m.Payload = payload

	return true
}

func (m *Message) ToTcpMessage() *messages.Message {
	msg := &messages.Message{
		Type:   "PUB",
		MsgId:  m.Id,
		Fields: make(map[string][]byte),
	}

	msg.Fields["msgId"] = []byte(m.Id)
	msg.Fields["routes"] = []byte(strings.Join(m.Routes, ","))
	msg.Fields["payload"] = m.Payload

	return msg
}

func (m *Message) ToBinary() ([]byte, error) {
	b, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (m *Message) FromBinary(b []byte) error {
	err := json.Unmarshal(b, m)
	if err != nil {
		return err
	}
	return nil
}
