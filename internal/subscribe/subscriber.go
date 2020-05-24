package subscribe

import (
	log "github.com/sirupsen/logrus"
	"go-broker/internal/serializer"
	"go-broker/internal/tcp"
	"sync"
	"time"
)

type Subscriber struct {
	clientId           string
	server             *tcp.Server
	timer              time.Timer
	sendMessageMap     map[string]*PublishedMessage
	queue              []*PublishedMessage
	concurrentMsgCount int
	mutex              sync.Mutex
}

func (s *Subscriber) enqueueMessage(message *PublishedMessage) {
	s.queue = append(s.queue, message)
	s.sendPendingMessages()
}

func (s *Subscriber) onMessageAck(msgId string) {
	msg := s.sendMessageMap[msgId]

	if msg != nil {
		delete(s.sendMessageMap, msgId)
		s.concurrentMsgCount -= 1
	}
}

func (s *Subscriber) onMessageNack(msgId string) {
	msg := s.sendMessageMap[msgId]

	if msg != nil {
		delete(s.sendMessageMap, msgId)
		s.concurrentMsgCount -= 1
	}
}

func (s *Subscriber) sendPendingMessages() {
	for {
		if s.concurrentMsgCount <= 0 {
			continue
		}
		if len(s.queue) == 0 {
			return
		}
		s.concurrentMsgCount += 1

		msg := s.queue[0]
		s.queue = s.queue[1:]

		s.sendMessageMap[msg.MsgId] = msg

		ser := serializer.NewLineSeparatedSerializer()

		ser.WriteStr("msgId", msg.MsgId)
		ser.WriteBytes("payload", msg.Payload)

		err := s.server.SendToClient(s.clientId, ser.GetMessageBytes())

		if err != nil {
			log.Errorf("error while sending to subscriber, err: %s", err)
		}

		log.Infof("Subscriber, sent msg with id: %s to client %s", msg.MsgId, s.clientId)

	}
}
