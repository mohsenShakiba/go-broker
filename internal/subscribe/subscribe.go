package subscribe

import "go-broker/internal/tcp"

type SubscriberManager struct {
	ackMessageChan     chan<- string
	nackMessageChan    chan<- string
	publishMessageChan <-chan PublishedMessage
	clients            map[string][]string
}

type PublishedMessage struct {
	Payload []byte
	Routes  []string
}

func InitSubscriberManager(ackChan chan<- string, nackChan chan<- string, socketServer *tcp.Server, publishMessageChan <-chan PublishedMessage) {
	mgr := SubscriberManager{
		ackMessageChan:     ackChan,
		nackMessageChan:    nackChan,
		publishMessageChan: publishMessageChan,
		clients:            make(map[string][]string),
	}

	socketServer.RegisterHandler("SUB", mgr.handleSubscribeMessage)
	socketServer.RegisterHandler("ACK", mgr.handleAckMessage)
	socketServer.RegisterHandler("NCK", mgr.handleNackMessage)

	go func() {
		for {
			msg := <-mgr.publishMessageChan

		}
	}()
}

func (s *SubscriberManager) handleSubscribeMessage(msgContext *tcp.MessageContext) {

}

func (s *SubscriberManager) handleAckMessage(msgContext *tcp.MessageContext) {

}

func (s *SubscriberManager) handleNackMessage(msgContext *tcp.MessageContext) {

}

func (s *SubscriberManager) processMessage(msg PublishedMessage) {
	// check which client should receive the message
	// enqueue the message in the client
}
