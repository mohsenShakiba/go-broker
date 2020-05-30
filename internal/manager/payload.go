package manager

import (
	"go-broker/internal/tcp/messages"
	"strings"
)

// PayloadMessage message contains the actual payload data and routing data
type PayloadMessage struct {
	Id      string
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

	p.Id = msg.MsgId
	p.Routes = routes
	p.Payload = payload

	return true
}

func (p *PayloadMessage) ToTcpMessage() *messages.Message {
	msg := &messages.Message{
		Type:   "PUB",
		MsgId:  p.Id,
		Fields: make(map[string][]byte),
	}

	msg.Fields["msgId"] = []byte(p.Id)
	msg.Fields["routes"] = []byte(strings.Join(p.Routes, ","))
	msg.Fields["payload"] = p.Payload

	return msg
}
