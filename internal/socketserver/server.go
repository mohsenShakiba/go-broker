package socketserver

import (
	"fmt"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"net"
	"time"
)

type Server struct {
	config SocketServerConfig
	listener net.Listener
	messageChan chan *clientMessage
	clients []*SocketClient
	publishedMessageChan chan<- string
}

func  Init(config SocketServerConfig, publishMessageChan chan<- string) Server {

	s := Server{}

	s.messageChan = make(chan *clientMessage, 1000)
	s.clients = make([]*SocketClient, 0, 100)
	s.publishedMessageChan = publishMessageChan
	s.config = config

	go s.listen()

	return s
}

func (s *Server) listen() {


	address := fmt.Sprintf(":%d", s.config.ConnectionPort)

	ln, err := net.Listen("tcp", address)

	if err != nil {
		log.Fatalf("couldn't listen on specified port, err: %s", err)
	}

	log.Printf("started listening on port %d", s.config.ConnectionPort)

	defer ln.Close()

	go s.listenToClientEvents()

	for {

		c, err := ln.Accept()

		log.Infof("a socket connection was established from %s", c.RemoteAddr())

		if err != nil {
			log.Infof("error accored while accepting connection, err: %s", err)
			return
		}

		go s.handleConnection(c)
	}

}

func (s *Server) listenToClientEvents() {
	msg := <-s.messageChan

	switch msg.Type {
	case clientMessageTypeDisconnect:
		s.removeClient(msg.ClientId)
		break
	case clientMessageTypePublish:
		s.publishedMessageChan <- msg.Payload
		break
	}
}

func (s *Server) handleConnection(conn net.Conn) {

	clientId := uuid.New()

	client := &SocketClient{
		clientId:        clientId.String(),
		connectionEpoch: time.Now().Unix(),
		clientType:      ClientUndetermined,
		isAuthenticated: false,
		connection:      conn,
		onMessageChan:   s.messageChan,
	}

	log.Infof("added a new client with id: %s", clientId)

	s.clients = append(s.clients, client)

	credStore := credentialStore{
		config: s.config,
	}

	client.startAuthenticate(credStore)

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



