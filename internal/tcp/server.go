package tcp

import (
	"fmt"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"go-broker/internal/message"
	"net"
)

// ServerConfig will hold the information about the behavior of the socket server
type ServerConfig struct {
	ConnectionPort int32
}

type ClientMessage struct {
	Message *message.Message
	Client  *Client
}

// Server will start a TCP socket server to accept incoming connections and read the data from the conn
// the data is then sent to a chanel which is processed by the manager
type Server struct {
	config   ServerConfig
	listener net.Listener
	handlers map[string]func(msgHandler *ClientMessage)
}

// Init will create a new socket server
func Init(config ServerConfig) *Server {

	s := &Server{
		config:   config,
		handlers: make(map[string]func(msgHandler *ClientMessage)),
	}

	return s
}

// Start will start the socket server
func (s *Server) Start() {

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", s.config.ConnectionPort))

	if err != nil {
		log.Fatalf("couldn't listen on specified port, err: %s", err)
	}

	defer listener.Close()

	log.Infof("started listening on port %d", s.config.ConnectionPort)

	for {

		c, err := listener.Accept()

		if err != nil {
			log.Errorf("error while accepting conn, err: %s", err)
			continue
		}

		go s.handleConnection(c)
	}

}

func (s *Server) RegisterHandler(t string, h func(msg *ClientMessage)) {
	s.handlers[t] = h
}

// handleConnection will start accepting connections
func (s *Server) handleConnection(conn net.Conn) {

	clientId := uuid.New()

	client := &Client{
		ClientId: clientId.String(),
		conn:     conn,
	}

	log.Infof("added a new client with Id: %s", clientId)

	for {
		msg, ok := client.Read()

		if !ok {
			return
		}

		s.processMessage(client, msg)
	}

}

func (s *Server) processMessage(client *Client, msg *message.Message) {

	handler, ok := s.handlers[msg.Type]

	if !ok {
		log.Errorf("no handler was found for %s", msg.Type)
	}

	go handler(&ClientMessage{
		Message: msg,
		Client:  client,
	})

}
