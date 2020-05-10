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
	messageChan          chan clientMessage
	clients              []*socketClient
	publishedMessageChan chan<- string
}

func Init(config SocketServerConfig, publishMessageChan chan<- string) *Server {

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
		connectionEpoch: time.Now().Unix(),
		clientType:      clientUndetermined,
		isClosed:        false,
		isAuthenticated: false,
		connection:      conn,
		onMessageChan:   s.messageChan,
	}

	log.Infof("added a new client with id: %s", clientId)

	s.clients = append(s.clients, client)

	client.startReceive()

}

// this method will listen to event from socket clients
func (s *Server) listenToClientEvents() {
	clientMsg := <-s.messageChan

	// get client
	client := s.findClientById(clientMsg.clientId)

	if client == nil {
		log.Errorf("a message was published from unknown client with id: %s", clientMsg.clientId)
	}

	// parse the message
	parsedMsg := parseMessage(clientMsg.payload)

	// detect the type of message
	switch v := parsedMsg.(type) {
	case authenticateMessage:
		// authenticate the client
		break
	case routedMessage:
		break
	case subscribeMessage:
		break
	case ackMessage:
		break
	case nackMessage:
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

func (s *Server) authenticateClient(client *socketClient, msg authenticateMessage) {

	if client.isAuthenticated {
		log.Warnf("the client %s has already been authenticated, ignoring")
		return
	}

	cred := fmt.Sprintf("%s:%s", msg.userName, msg.password)

	for _, validCred := range s.config.Credentials {
		if cred == validCred {
			// set a authenticated
			client.setAsAuthenticated()

			// send authentication event
			recEv := receiveMessage{
				id: msg.id,
			}

			// send event
			client.send(recEv.format())
		}
	}

	// close the client
	client.close()

}
