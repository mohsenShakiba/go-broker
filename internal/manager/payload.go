package manager

import (
	"go-broker/internal/tcp/messages"
	"strings"
)

// PayloadMessage message contains the actual payload data and routing data
type PayloadMessage struct {
	MsgId   string
	Routes  []string
	Payload []byte
}

func (p *PayloadMessage) FromTcpMessage(msg *messages.Message) bool {

	routesStr, ok := msg.ReadStr("routes")

	if !ok {
		return false
	}

	routes := strings.Split(routesStr, ",")

	payload, ok := msg.ReadByteArr("payload")

	if !ok {
		return false
	}

	p.MsgId = msg.MsgId
	p.Routes = routes
	p.Payload = payload

	return true
}

func (p *PayloadMessage) ToTcpMessage() *messages.Message {
	msg := &messages.Message{
		Type:   "PUB",
		MsgId:  p.MsgId,
		Fields: make(map[string][]byte),
	}

	msg.Fields["msgId"] = []byte(p.MsgId)
	msg.Fields["routes"] = []byte(strings.Join(p.Routes, ","))
	msg.Fields["payload"] = p.Payload

	return msg
}
