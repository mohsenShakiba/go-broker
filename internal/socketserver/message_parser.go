package socketserver

import (
	"errors"
	"fmt"
	"go-broker/internal/socketserver/serializer"
)

func parseMessage(message []byte) (interface{}, error) {

	// extract the message type
	msgType := message[:3]

	payload := message[3:]

	s := serializer.NewJsonSerializer()

	// parse based on type
	switch string(msgType) {
	case authenticateMessageType:
		msg := &authenticateMessage{}
		err := s.Deserialize(payload, msg)
		return msg, err
	case routedMessageType:
		msg := &routedMessage{}
		err := s.Deserialize(payload, msg)
		return msg, err
	case subscribeMessageType:
		msg := &subscribeMessage{}
		err := s.Deserialize(payload, msg)
		return msg, err
	case ackMessageType:
		msg := &ackMessage{}
		err := s.Deserialize(payload, msg)
		return msg, err
	case nackMessageType:
		msg := &nackMessage{}
		err := s.Deserialize(payload, msg)
		return msg, err

	}

	errMsg := fmt.Sprintf("the message type %s cannot be parsed", string(msgType))
	return nil, errors.New(errMsg)

}
