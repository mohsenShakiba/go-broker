package socketserver

import (
	"fmt"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"net"
)

type SocketServerConfig struct {
	Credentials    []string
	ConnectionPort int32
}

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
		log.Errorf("a tcpMessage was published from unknown client with Id: %s", clientMsg.clientId)
	}

	// parse the tcpMessage
	msgContext := convertToMessage(clientMsg.payload)

}

func (s *Server) findClientById(id string) *socketClient {
	for _, c := range s.clients {
		if c.clientId == id {
			return c
		}
	}
	return nil
}
