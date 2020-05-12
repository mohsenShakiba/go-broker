package socketserver

import (
	"fmt"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"go-broker/internal/socketserver/serializer"
	"net"
)

type Server struct {
	config               SocketServerConfig
	listener             net.Listener
	messageChan          chan clientMessage
	clients              []*socketClient
	publishedMessageChan chan<- ServerEvents
}

func Init(config SocketServerConfig, publishMessageChan chan<- ServerEvents) *Server {

	s := &Server{}

	s.messageChan = make(chan clientMessage, 100)
	s.clients = make([]*socketClient, 0, 100)
	s.publishedMessageChan = publishMessageChan
	s.config = config

	go s.listen()

	go s.listenToClientEvents()

	return s
}

// this method will fireup the socket server
// it will also accept connection
func (s *Server) listen() {

	address := fmt.Sprintf(":%d", s.config.ConnectionPort)

	ln, err := net.Listen("tcp", address)

	if err != nil {
		log.Fatalf("couldn't listen on specified port, err: %s", err)
	}

	log.Infof("started listening on port %d", s.config.ConnectionPort)

	defer ln.Close()

	for {

		c, err := ln.Accept()

		log.Infof("a socket connection was established from %s", c.RemoteAddr())

		if err != nil {
			log.Errorf("error while accepting connection, err: %s", err)
			return
		}

		go s.handleConnection(c)
	}

}

func (s *Server) handleConnection(conn net.Conn) {

	clientId := uuid.New()

	client := &socketClient{
		clientId:        clientId.String(),
		clientType:      clientUndetermined,
		isClosed:        false,
		isAuthenticated: false,
		connection:      conn,
		onMessageChan:   s.messageChan,
	}

	log.Infof("added a new client with Id: %s", clientId)

	s.clients = append(s.clients, client)

	client.startReceive()

}

func (s *Server) listenToClientEvents() {
	for {
		clientMsg := <-s.messageChan
		s.processClientEvents(clientMsg)
	}
}

// this method will listen to event from socket clients
func (s *Server) processClientEvents(clientMsg clientMessage) {

	// get client
	client := s.findClientById(clientMsg.clientId)

	if client == nil {
		log.Errorf("a message was published from unknown client with Id: %s", clientMsg.clientId)
	}

	// parse the message
	parsedMsg, err := parseMessage(clientMsg.payload)

	if err != nil {
		log.Errorf("parsing message failed, error: %s", err)
	}

	// detect the type of message
	switch v := parsedMsg.(type) {
	case *authenticateMessage:
		s.authenticateClient(client, v)
		break
	case *routedMessage:
		s.processRoutedMessage(client, v)
		break
	case *subscribeMessage:
		s.processSubscribeMessage(client, v)
		break
	case *ackMessage:
		s.processAckMessage(client, v)
		break
	case *nackMessage:
		s.processNackMessage(client, v)
		break
	}

}

func (s *Server) findClientById(id string) *socketClient {
	for _, c := range s.clients {
		if c.clientId == id {
			return c
		}
	}
	return nil
}

func (s *Server) authenticateClient(client *socketClient, msg *authenticateMessage) {

	if client.isAuthenticated {
		log.Warnf("the client %s has already been authenticated, ignoring", client.clientId)
		return
	}

	cred := fmt.Sprintf("%s:%s", msg.UserName, msg.Password)
	binarySerializer := serializer.NewJsonSerializer()

	for _, validCred := range s.config.Credentials {
		if cred == validCred {
			// set a authenticated
			client.setAsAuthenticated()

			// send authentication event
			successEv := receiveMessage{
				Id:      msg.Id,
				Success: true,
			}

			successEvBinary, err := binarySerializer.Serialize(&successEv)

			if err != nil {
				log.Errorf("failed to serialize authenticate event, error: %s", err)
				break
			}

			// send event
			log.Infof("authentication was successful for client %s", client.clientId)
			client.send(successEvBinary)

			return
		}
	}

	// send authentication failed event
	failedEv := receiveMessage{
		Id:      msg.Id,
		Success: false,
	}

	log.Infof("authentication failed for client %s", client.clientId)
	failedEvBinary, _ := binarySerializer.Serialize(&failedEv)

	// send event
	client.send(failedEvBinary)

}

func (s *Server) processRoutedMessage(client *socketClient, msg *routedMessage) {

	// check if client is authenticated
	if !client.isAuthenticated {
		log.Warnf("the client %s isn't authenticated to send routed messages, ignoring", client.clientId)
		return
	}

	// send message to manager
	ev := &PublishMessageEvent{
		ClientId: client.clientId,
		MsgId:    msg.Id,
		Routes:   msg.Routes,
		Payload:  msg.Payload,
	}

	s.publishedMessageChan <- ev

}

func (s *Server) processSubscribeMessage(client *socketClient, msg *subscribeMessage) {

	// check if client is authenticated
	if !client.isAuthenticated {
		log.Warnf("the client %s isn't authenticated for subscription, ignoring", client.clientId)
		return
	}

	client.clientType = clientSubscriber

	log.Infof("the client %s was registered as subscriber", client.clientId)

	// send subscription config to sender
	ev := &SubscriberAuthenticatedEvent{
		ClientId: client.clientId,
		Routes:   msg.Routes,
		BufSize:  msg.BufSize,
	}

	s.publishedMessageChan <- ev

	// send the rec event
	recMsg := receiveMessage{
		Id:      msg.Id,
		Success: true,
	}

	binarySerializer := serializer.NewJsonSerializer()

	recMsgSerialized, _ := binarySerializer.Serialize(&recMsg)

	client.send(recMsgSerialized)
}

func (s *Server) processAckMessage(client *socketClient, msg *ackMessage) {

	// check if client is authenticated
	if !client.isAuthenticated {
		log.Warnf("the client %s isn't authenticated for ack, ignoring", client.clientId)
		return
	}

	// send ack to manager

}

func (s *Server) processNackMessage(client *socketClient, msg *nackMessage) {

	// check if client is authenticated
	if !client.isAuthenticated {
		log.Warnf("the client %s isn't authenticated for nack, ignoring", client.clientId)
		return
	}

	// send nack to manager

}
