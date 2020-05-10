package socketserver

import (
	"net"
)

type socketClient struct {
	clientId        string
	connectionEpoch int64
	isAuthenticated bool
	clientType      int
	isClosed        bool
	connection      net.Conn
	onMessageChan   chan<- clientMessage
}

type clientMessage struct {
	clientId   string
	clientType int
	payload    []byte
}

const (
	clientUndetermined = 0
	clientPublisher    = 1
	clientSubscriber   = 2
)

//// this method will read the handshake request from client
//// handshake is required to perform authentication
//func (c *socketClient) initHandshake(store credentialStore) {
//
//	log.Infof("waiting for client with id: %s to authenticate", c.clientId)
//
//	// read the first message from socket
//	// the first message will contain the handshake information
//	handshakeRequest, err := c.receive()
//
//	if err != nil {
//		c.close()
//		return
//	}
//
//	// the handshake request must at least be 2 bytes
//	if len(handshakeRequest) <= 1 {
//		log.Errorf("handshake request was too small")
//		c.close()
//		return
//	}
//
//	// split the handshake request
//	clientType := handshakeRequest[:1]
//	credential := handshakeRequest[1:]
//
//	isAuthenticated := c.authenticate(store, credential)
//
//	if !isAuthenticated {
//		c.close()
//		return
//	}
//
//	typeResult, ok := c.parseClientType(clientType)
//
//	if !ok {
//		c.close()
//		return
//	}
//
//	// set authentication and client type
//	c.isAuthenticated = true
//	c.clientType = typeResult
//
//	err = c.send([]byte("1"))
//
//	if err != nil {
//		log.Errorf("could not sent authentication success event, error: %s", err)
//		c.close()
//		return
//	}
//
//	go c.startReceive()
//}
//
//func (c *socketClient) authenticate(store credentialStore, cred string) bool {
//
//	if store.isValid(cred) {
//		log.Infof("credential submitted: %s is valid", cred)
//		return true
//	} else {
//		log.Warnf("credential submitted: %s is invalid", cred)
//		return false
//	}
//}
//
//func (c *socketClient) parseClientType(typeStr string) (int, bool) {
//	ctype, _ := strconv.Atoi(typeStr)
//
//	switch ctype {
//	case clientPublisher:
//		return ctype, true
//	case clientSubscriber:
//		return ctype, true
//	}
//
//	log.Errorf("the client didn't provide a valid type, type: %s", typeStr)
//	return 0, false
//}

func (c *socketClient) setAsAuthenticated() {
	c.isAuthenticated = true
}

func (c *socketClient) setClientType(clientType int) {
	c.clientType = clientType
}

func (c *socketClient) startReceive() {
	for {

		if c.isClosed {
			return
		}

		msg, err := read(c.connection, 1024)

		if err != nil {
			continue
		}

		c.onMessageChan <- msg

	}
}

func (c *socketClient) close() {
	_ = c.connection.Close()
	c.isClosed = true
}

func (c *socketClient) send(payload []byte) error {
	return write(c.connection, payload)
}
