package manager

import (
	"encoding/json"
	"go-broker/internal/tcp/messages"
	"hash/fnv"
	"strings"
)

// PayloadMessage message contains the actual payload data and routing data
type PayloadMessage struct {
	Id      string
	Routes  []string
	Payload []byte
}

func getStringHash(id string) int64 {
	h := fnv.New32a()
	h.Write([]byte(id))
	return int64(h.Sum32())
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

func (p *PayloadMessage) ToBinary() ([]byte, error) {
	b, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (p *PayloadMessage) FromBinary(b []byte) error {
	err := json.Unmarshal(b, p)
	if err != nil {
		return err
	}
	return nil
}
