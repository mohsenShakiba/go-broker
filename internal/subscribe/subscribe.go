package subscribe

import "go-broker/internal/tcp"

type SubscriberManager struct {
	ackMessageChan     chan<- string
	nackMessageChan    chan<- string
	publishMessageChan <-chan PublishMessage
}

type PublishMessage struct {
	Routes  []string
	Payload []byte
	MsgId   string
}

func InitSubscriberManager(ackChan chan<- string, nackChan chan<- string, socketServer *tcp.Server) {
	mgr := SubscriberManager{
		ackMessageChan:  ackChan,
		nackMessageChan: nackChan,
	}

	socketServer.RegisterHandler("SUB", mgr.handleSubscribeMessage)
	socketServer.RegisterHandler("ACK", mgr.handleAckMessage)
	socketServer.RegisterHandler("NCK", mgr.handleNackMessage)

}

func (s *SubscriberManager) handleSubscribeMessage(msgContext *tcp.MessageContext) {

}

func (s *SubscriberManager) handleAckMessage(msgContext *tcp.MessageContext) {

}

func (s *SubscriberManager) handleNackMessage(msgContext *tcp.MessageContext) {

}
