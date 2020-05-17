package publish

import (
	log "github.com/sirupsen/logrus"
	"go-broker/internal/tcp"
	"strings"
)

type PublishedMessage struct {
	Routes  []string
	Payload []byte
}

type PublisherManager struct {
	publishedMessageChan chan<- *PublishedMessage
}

func InitPublisherManager(ch chan<- *PublishedMessage, server *tcp.Server) *PublisherManager {

	receiver := &PublisherManager{
		publishedMessageChan: ch,
	}

	// add handlers
	server.RegisterHandler("PUb", receiver.handlePublishMessage)

	return receiver
}

func (p *PublisherManager) handlePublishMessage(msgContext *tcp.MessageContext) {

	// get message routes
	msgId, ok := msgContext.GetMessageId()

	if !ok {
		log.Errorf("the published message doesn't have a valid msgId, discarding message")
		return
	}

	// get message routes
	routesStr, ok := msgContext.ReadStr("routes")

	if !ok {
		log.Errorf("the published message doesn't have a valid route, discarding message")
		return
	}

	// get message payload
	payload, ok := msgContext.ReadByteArr("payload")

	if !ok {
		log.Errorf("the published message doesn't have a valid payload, discarding message")
		return
	}

	// parse message routes
	routesArr := strings.Split(routesStr, ",")

	// create published message
	msg := &PublishedMessage{
		Routes:  routesArr,
		Payload: payload,
	}

	p.publishedMessageChan <- msg

	err := msgContext.SendAck()

	if err != nil {
		log.Errorf("failed to send message ack for id: %s", msgId)
	}
}
