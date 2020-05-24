package tcp

import (
	"fmt"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"net"
)

// ServerConfig will hold the information about the behavior of the socket server
type ServerConfig struct {
	ConnectionPort int32
}

type Message struct {
	Payload []byte
	Client  *Client
}

// Server will start a TCP socket server to accept incoming connections and read the data from the conn
// the data is then sent to a chanel which is processed by the manager
type Server struct {
	config      ServerConfig
	listener    net.Listener
	clients     []*Client
	messageChan chan<- Message
	//handlers      map[string]func(msgHandler *MessageContext)
}

// Init will create a new socket server
func Init(config ServerConfig, messageChan chan<- Message) *Server {

	s := &Server{
		config:      config,
		listener:    nil,
		clients:     make([]*Client, 0, 100),
		messageChan: messageChan,
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

// handleConnection will start accepting connections
func (s *Server) handleConnection(conn net.Conn) {

	clientId := uuid.New()

	client := &Client{
		ClientId: clientId.String(),
		conn:     conn,
	}

	log.Infof("added a new client with Id: %s", clientId)

	s.clients = append(s.clients, client)

}
