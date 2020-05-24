package publish

import (
	log "github.com/sirupsen/logrus"
	"go-broker/internal/message"
	"go-broker/internal/tcp"
	"strings"
)

type PublisherManager struct {
	publishedMessageChan chan<- *message.Message
}

func InitPublisherManager(ch chan<- *message.Message, server *tcp.Server) *PublisherManager {

	receiver := &PublisherManager{
		publishedMessageChan: ch,
	}

	// add handlers
	server.RegisterHandler("PUB", receiver.handlePublishMessage)

	return receiver
}

func (p *PublisherManager) handlePublishMessage(clientMessage *tcp.ClientMessage) {

	// get message routes
	msgId, ok := clientMessage.Message.ReadMsgId()

	if !ok {
		log.Errorf("the published message doesn't have a valid msgId, discarding message")
		return
	}

	// get message routes
	routesStr, ok := clientMessage.Message.ReadStr("routes")

	if !ok {
		log.Errorf("the published message doesn't have a valid route, discarding message")
		return
	}

	// get message payload
	payload, ok := clientMessage.Message.ReadByteArr("payload")

	if !ok {
		log.Errorf("the published message doesn't have a valid payload, discarding message")
		return
	}

	// parse message routes
	routesArr := strings.Split(routesStr, ",")

	log.Infof("sending published message: %s to manager", msgId)

	p.publishedMessageChan <- msg

	err := msgContext.SendAck()

	if err != nil {
		log.Errorf("failed to send message ack for id: %s", msgId)
	}
}
