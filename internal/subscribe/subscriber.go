package subscribe

import (
	"go-broker/internal/tcp"
	"go-broker/internal/tcp/serializer"
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
		s.mutex.Lock()
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
		s.mutex.Unlock()

		ser := serializer.NewLineSeparatedSerializer()

		ser.WriteStr("msgId", msg.MsgId)
		ser.WriteBytes("payload", msg.Payload)

		s.server.SendToClient(s.clientId, []byte(ser.GetMessagePrefix()))

		for _, b := range ser.Bytes {
			s.server.SendToClient(s.clientId, b)
		}

	}
}
