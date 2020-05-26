package publish

import (
	log "github.com/sirupsen/logrus"
	"go-broker/internal/manager/messages"
	"go-broker/internal/tcp"
	"strings"
)

type PublisherManager struct {
	publishedMessageChan chan<- *messages.Message
}

func InitPublisherManager(ch chan<- *messages.Message, server *tcp.Server) *PublisherManager {

	receiver := &PublisherManager{
		publishedMessageChan: ch,
	}

	// add handlers
	server.RegisterHandler("PUB", receiver.handlePublishMessage)

	return receiver
}

func (p *PublisherManager) handlePublishMessage(clientMessage *tcp.Context) {

	// get messages routes
	msgId, ok := clientMessage.Message.ReadMsgId()

	if !ok {
		log.Errorf("the published messages doesn't have a valid msgId, discarding messages")
		return
	}

	// get messages routes
	routesStr, ok := clientMessage.Message.ReadStr("routes")

	if !ok {
		log.Errorf("the published messages doesn't have a valid route, discarding messages")
		return
	}

	// get messages payload
	payload, ok := clientMessage.Message.ReadByteArr("payload")

	if !ok {
		log.Errorf("the published messages doesn't have a valid payload, discarding messages")
		return
	}

	// parse messages routes
	routesArr := strings.Split(routesStr, ",")

	log.Infof("sending published messages: %s to manager", msgId)

	p.publishedMessageChan <- msg

	err := msgContext.SendAck()

	if err != nil {
		log.Errorf("failed to send messages ack for id: %s", msgId)
	}
}
