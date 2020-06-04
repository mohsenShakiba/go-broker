package tcp

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net"
)

// ServerConfig will hold the information about the behavior of the socket server
type ServerConfig struct {
	ConnectionPort int32
}

// Server will start a TCP socket server to accept incoming connections and read the data from the conn
// the data is then sent to a chanel which is processed by the manager
type Server struct {
	config   ServerConfig
	listener net.Listener
	*Multiplexer
}

// New will create a new socket server
func New(config ServerConfig) *Server {

	s := &Server{
		config:      config,
		Multiplexer: newMultiplexer(),
	}

	return s
}

// Start will start the socket server
func (s *Server) Start() {

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", s.config.ConnectionPort))

	if err != nil {
		log.Fatalf("couldn't listen on specified port, err: %s", err)
	}

	log.Infof("started listening on port %d", s.config.ConnectionPort)

	go func() {

		defer listener.Close()

		for {

			c, err := listener.Accept()

			if err != nil {
				log.Errorf("error while accepting conn, err: %s", err)
				continue
			}

			go s.handleConnection(c)
		}
	}()

}

// handleConnection will start accepting connections
func (s *Server) handleConnection(conn net.Conn) {

	client := initSocketClient(conn)

	log.Infof("added a new client with Id: %s", client.ClientId)

	for {
		msg, ok := client.Read()

		if !ok {
			return
		}

		s.process(client, msg)
	}

}
