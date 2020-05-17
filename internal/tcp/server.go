package tcp

import (
	"fmt"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"net"
)

type TcpConfig struct {
	Credentials    []string
	ConnectionPort int32
}

type Server struct {
	config        TcpConfig
	listener      net.Listener
	clients       []*socketClient
	clientMsgChan chan clientMessage
	handlers      map[string]func(msgHandler *MessageContext)
}

func Init(config TcpConfig, publishMessageChan chan<- clientMessage) *Server {

	s := &Server{
		config:        config,
		listener:      nil,
		clients:       make([]*socketClient, 0, 100),
		clientMsgChan: make(chan clientMessage, 100),
		handlers:      make(map[string]func(msgHandler *MessageContext), 0),
	}

	go s.listen()

	go s.listenToClientEvents()

	return s
}

func (s *Server) RegisterHandler(t string, handler func(msgHandler *MessageContext)) {
	s.handlers[t] = handler
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
		clientId:      clientId.String(),
		isClosed:      false,
		connection:    conn,
		onMessageChan: s.clientMsgChan,
	}

	log.Infof("added a new client with Id: %s", clientId)

	s.clients = append(s.clients, client)

	client.startReceive()

}

func (s *Server) listenToClientEvents() {
	for {
		clientMsg := <-s.clientMsgChan
		s.processClientEvents(clientMsg)
	}
}

// this method will listen to event from socket clients
func (s *Server) processClientEvents(clientMsg clientMessage) {

	// get client
	client := s.findClientById(clientMsg.clientId)

	if client == nil {
		log.Errorf("a tcpMessage was published from unknown client with Id: %s", clientMsg.clientId)
	}

	// parse the tcpMessage
	msgContext := convertToMessage(clientMsg.payload, client)

	msgType, ok := msgContext.GetMessageType()

	if !ok {
		log.Errorf("the message doesn't seem to have a valid message type, msg type: %s", msgType)
		return
	}

	mh := s.handlers[msgType]

	if mh == nil {
		log.Errorf("no handler has been registered for type %s", msgType)
		return
	}

	mh(msgContext)

}

func (s *Server) findClientById(id string) *socketClient {
	for _, c := range s.clients {
		if c.clientId == id {
			return c
		}
	}
	return nil
}
