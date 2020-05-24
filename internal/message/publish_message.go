package message

import "strings"

type PublishMessage struct {
	MsgId   string
	Routes  []string
	Payload []byte
}

func (m *PublishMessage) FromMessage(msg *Message) bool {

	msgId, ok := msg.ReadMsgId()

	if !ok {
		return false
	}

	routesStr, ok := msg.ReadStr("routes")

	if !ok {
		return false
	}

	routes := strings.Split(routesStr, ",")

	payload, ok := msg.ReadByteArr("payload")

	if !ok {
		return false
	}

	m.MsgId = msgId
	m.Routes = routes
	m.Payload = payload

	return true
}

func (m *PublishMessage) ToMessage() *Message {
	msg := &Message{
		Type:   "PUB",
		Fields: make(map[string][]byte),
	}

	msg.Fields["msgId"] = []byte(m.MsgId)
	msg.Fields["routes"] = []byte(strings.Join(m.Routes, ","))
	msg.Fields["payload"] = m.Payload

	return msg
}
