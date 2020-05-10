package socketserver

import (
	"bytes"
	log "github.com/sirupsen/logrus"
)

func parseMessage(message []byte) baseMessage {

	// split the message based on the \n
	msgParts := bytes.Split(message, []byte("\n"))

	// check if length is more than 1
	if len(msgParts) <= 1 {
		log.Errorf("the message isn't in valid format, it must contain more than once part, message: %s", string(message))
		return nil
	}

	// extract the message type
	msgType := msgParts[0]

	payload := msgParts[1:]

	// parse based on type
	switch string(msgType) {
	case authenticateMessageType:
		return authenticateMessageFromPayload(payload)
	case routedMessageType:
		return routedMessageFromPayload(payload)
	case subscribeMessageType:
		return subscribeMessageFromPayload(payload)
	case ackMessageType:
		return ackMessageFromPayload(payload)
	case nackMessageType:
		return nackMessageFromPayload(payload)

	}

	return nil

}
