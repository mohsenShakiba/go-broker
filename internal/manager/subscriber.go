package manager

import (
	log "github.com/sirupsen/logrus"
	"go-broker/internal/message"
	"go-broker/internal/serializer"
	"go-broker/internal/subscribe"
	"go-broker/internal/tcp"
	"sync"
	"time"
)

type subscriberConfig struct {
	maxConcurrentMessageCount int
	routes                    []string
}

type Subscriber struct {
	client           tcp.Client
	config           subscriberConfig
	sentMessages     map[string]*message.PublishMessage
	sentMessageCount int
	queue            []*message.PublishMessage
	mutex            sync.Mutex
}

// start will start sending messages from the queue
func (s *Subscriber) start() {

	for {

		// sleep if no messages are in the queue
		if len(s.queue) == 0 {
			time.Sleep(time.Millisecond * 100)
		}

		// if sentMessageCount is equal to config maxConcurrentMessageCount
		if s.sentMessageCount >= s.config.maxConcurrentMessageCount {
			time.Sleep(time.Millisecond * 100)
		}

		// retrieve message
		msg := s.queue[0]
		s.queue = s.queue[1:]

		// add to sent messages
		s.sentMessages[msg.MsgId] = msg

		// send message

	}
}

func (s *Subscriber) enqueueMessage(message *subscribe.PublishedMessage) {
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
