package subscribe

import (
	"errors"
	"go-broker/internal/tcp"
	"time"
)

type Subscriber struct {
	clientId        string
	timer           time.Timer
	lastSentMessage *PublishMessage
	queue           chan *PublishMessage
}

const (
	queueFullError = "client queue if full"
)

func (s *Subscriber) beginProcess(server *tcp.Server) {
	for {
		msg := <-s.queue

		err := server.SendToClient(s.clientId, msg.Payload)

		if err != nil {
			return
		}
	}
}

func (s *Subscriber) enqueueMessage(message *PublishMessage) error {
	select {
	case s.queue <- message:
		return nil
	default:
		return errors.New(queueFullError)
	}
}

func (s *Subscriber) onMessageAck(msgId string) {

}

func (s *Subscriber) onMessageNack(msgId string) {

}

func (s *Subscriber) onMessageTimeout() {

}
