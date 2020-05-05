package socketserver

import (
	"fmt"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"net"
	"time"
)

type Server struct {
	config               SocketServerConfig
	listener             net.Listener
	messageChan          chan *clientMessage
	clients              []*SocketClient
	publishedMessageChan chan<- string
}

func Init(config SocketServerConfig, publishMessageChan chan<- string) *Server {

	s := &Server{}

	s.messageChan = make(chan *clientMessage, 1000)
	s.clients = make([]*SocketClient, 0, 100)
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

// this method will listen to event from socket clients
func (s *Server) listenToClientEvents() {
	msg := <-s.messageChan

	switch msg.Type {
	case clientMessageTypeDisconnect:
		s.removeClient(msg.ClientId)
	case clientMessageTypePublish:
		s.publishedMessageChan <- msg.Payload
	}
}

func (s *Server) handleConnection(conn net.Conn) {

	clientId := uuid.New()

	client := &SocketClient{
		clientId:        clientId.String(),
		connectionEpoch: time.Now().Unix(),
		clientType:      clientUndetermined,
		isAuthenticated: false,
		connection:      conn,
		onMessageChan:   s.messageChan,
	}

	log.Infof("added a new client with id: %s", clientId)

	s.clients = append(s.clients, client)

	credStore := credentialStore{
		config: s.config,
	}

	client.initHandshake(credStore)

}

func (s *Server) removeClient(clientId string) {
	index := -1

	for i, c := range s.clients {
		if c.clientId == clientId {
			index = i
		}
	}

	if index != -1 {
		s.clients[index] = s.clients[len(s.clients)-1]
		s.clients = s.clients[:len(s.clients)-1]
	}

}
