package socketserver

import (
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

const (
	authenticateMessageType = "AUT"
	routedMessageType       = "PUB"
	subscribeMessageType    = "SUB"
	ackMessageType          = "ACK"
	nackMessageType         = "NCK"
)

type clientPayloadMessage [][]byte

type baseMessage interface{}

// this message is sent to subscriber
type contentMessage struct {
	id      string
	payload []byte
}

type receiveMessage struct {
	id string
}

func (m *receiveMessage) format() []byte {
	return format([]byte("REC"), []byte(m.id))
}

// this message is sent by publisher and subscribers for authentication
type authenticateMessage struct {
	id       string
	userName string
	password string
}

// this message is sent by the publisher to a route
type routedMessage struct {
	id      string
	routes  []string
	payload []byte
}

// this message is sent by the subscriber
type subscribeMessage struct {
	id      string
	routes  []string
	bufSize int
}

// this message is sent by the subscriber to discard the message as processed
type ackMessage struct {
	id string
}

// this message is sent by the subscriber to requeue the message
type nackMessage struct {
	id string
}

func authenticateMessageFromPayload(payload clientPayloadMessage) authenticateMessage {

	if len(payload) != 3 {
		log.Errorf("the authentication payload isn't in valid format, payload: %s", payload)
	}

	return authenticateMessage{
		id:       string(payload[0]),
		userName: string(payload[1]),
		password: string(payload[2]),
	}
}

func routedMessageFromPayload(payload clientPayloadMessage) routedMessage {

	if len(payload) != 3 {
		log.Errorf("the routed payload isn't in valid format, payload: %s", payload)
	}

	return routedMessage{
		id:      string(payload[0]),
		routes:  strings.Split(string(payload[1]), ","),
		payload: payload[2],
	}
}

func subscribeMessageFromPayload(payload clientPayloadMessage) subscribeMessage {

	if len(payload) != 3 {
		log.Errorf("the subscription payload isn't in valid format, payload: %s", payload)
	}

	buffSize, err := strconv.Atoi(string(payload[2]))

	if err != nil {
		log.Errorf("the subscription doesn't provide a valid buffer size")
	}

	return subscribeMessage{
		id:      string(payload[0]),
		routes:  strings.Split(string(payload[1]), ","),
		bufSize: buffSize,
	}
}

func ackMessageFromPayload(payload clientPayloadMessage) ackMessage {

	if len(payload) != 1 {
		log.Errorf("the ack payload isn't in valid format, payload: %s", payload)
	}

	return ackMessage{
		id: string(payload[0]),
	}
}

func nackMessageFromPayload(payload clientPayloadMessage) nackMessage {

	if len(payload) != 1 {
		log.Errorf("the ack payload isn't in valid format, payload: %s", payload)
	}

	return nackMessage{
		id: string(payload[0]),
	}
}
