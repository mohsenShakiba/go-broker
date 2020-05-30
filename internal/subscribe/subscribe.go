package subscribe

import (
	log "github.com/sirupsen/logrus"
	"go-broker/internal/manager"
	"go-broker/internal/tcp"
	"strings"
	"sync"
	"time"
)

type SubscriberManager struct {
	messageResultChan  chan<- *MessageResult
	publishMessageChan <-chan *PublishedMessage
	routesMapping      map[string][]*manager.Subscriber
	clientMapping      map[string]*manager.Subscriber
	socketServer       *tcp.Server
}

const (
	MessageAck   = "msg_ack"
	MessageNack  = "msg_nack"
	NoSubscriber = "msg_no_subscriber"
)

type PublishedMessage struct {
	MsgId   string
	Payload []byte
	Routes  []string
}

type MessageResult struct {
	MsgId  string
	Result string
}

func InitSubscriberManager(socketServer *tcp.Server, publishMessageChan <-chan *PublishedMessage) *SubscriberManager {
	mgr := SubscriberManager{
		messageResultChan:  nil,
		publishMessageChan: publishMessageChan,
		routesMapping:      make(map[string][]*manager.Subscriber),
		clientMapping:      make(map[string]*manager.Subscriber),
		socketServer:       socketServer,
	}

	socketServer.RegisterHandler("SUB", mgr.handleSubscribeMessage)
	socketServer.RegisterHandler("ACK", mgr.handleAckMessage)
	socketServer.RegisterHandler("NCK", mgr.handleNackMessage)

	go mgr.processMessageQueue()

	return &mgr
}

func (s *SubscriberManager) handleSubscribeMessage(msgContext *tcp.MessageContext) {

	clientId := msgContext.ClientId

	routesStr, ok := msgContext.ReadStr("routes")

	if !ok {
		log.Errorf("could not read routes from subscriber client, discarding client")
		msgContext.Close()
	}

	routes := strings.Split(routesStr, ",")

	if len(routes) <= 0 {
		log.Errorf("the client didn't provide a valid route")
		msgContext.Close()
	}

	//send ack
	err := msgContext.SendAck()

	if err != nil {
		log.Errorf("error while sending ack for subscribe messages")
	}

	client := manager.Subscriber{
		clientId:           clientId,
		server:             s.socketServer,
		timer:              time.Timer{},
		sendMessageMap:     make(map[string]*PublishedMessage),
		queue:              make([]*PublishedMessage, 0),
		concurrentMsgCount: 1,
		mutex:              sync.Mutex{},
	}

	s.clientMapping[clientId] = &client

	log.Infof("added subscriber with client id: %s for routes: %s", clientId, routesStr)

	for _, route := range routes {

		clients := s.routesMapping[route]

		if clients == nil {
			clients = make([]*manager.Subscriber, 0)
		}

		clients = append(clients, &client)

		s.routesMapping[route] = clients

		log.Infof("client with id: %s was added to subscriber list", client.clientId)
	}
}

func (s *SubscriberManager) handleAckMessage(msgContext *tcp.MessageContext) {
	clientId := msgContext.ClientId

	msgId, ok := msgContext.GetMessageId()

	if !ok {
		log.Errorf("could not read msg id from ack messages, discarding")
		return
	}

	client := s.clientMapping[clientId]

	if client == nil {
		log.Errorf("a messages was received from a client that doesn't seem to exist, discarding")
		return
	}

	client.onMessageAck(msgId)

}

func (s *SubscriberManager) handleNackMessage(msgContext *tcp.MessageContext) {
	clientId := msgContext.ClientId

	msgId, ok := msgContext.GetMessageId()

	if !ok {
		log.Errorf("could not read msg id from nack messages, discarding")
		return
	}

	client := s.clientMapping[clientId]

	if client == nil {
		log.Errorf("a messages was received from a client that doesn't seem to exist, discarding")
		return
	}

	client.onMessageNack(msgId)
}

func (s *SubscriberManager) processMessageQueue() {

	func() {
		for {
			msg := <-s.publishMessageChan
			s.processMsg(msg)
		}
	}()
}

func (s *SubscriberManager) processMsg(msg *PublishedMessage) {
	// check which client should receive the messages
	// enqueue the messages in the client
	log.Infof("Submanager, messages received from publisher msgId: %s", msg.MsgId)
	subscriberFound := false
	for _, msgRoute := range msg.Routes {

		clients := s.routesMapping[msgRoute]

		if clients == nil {
			continue
		}

		if len(clients) == 0 {
			continue
		}

		for _, client := range clients {
			log.Infof("Submanager, enqueue messages with msgId: %s to client %s", msg.MsgId, client.clientId)
			client.enqueueMessage(msg)
			subscriberFound = true
		}

	}

	if !subscriberFound {
		s.setMsgResult(msg, NoSubscriber)
	}
}

func (s *SubscriberManager) setMsgResult(msg *PublishedMessage, result string) {
	//s.messageResultChan <- &MessageResult{
	//	Id:  msg.Id,
	//	Result: result,
	//}
}
