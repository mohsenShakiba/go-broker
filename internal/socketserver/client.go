package socketserver

import (
	log "github.com/sirupsen/logrus"
	"net"
	"strconv"
)

type SocketClient struct {
	clientId        string
	connectionEpoch int64
	clientType      int
	isAuthenticated bool
	connection      net.Conn
	onMessageChan   chan<- *clientMessage
}

const (
	ClientUndetermined = 0
	ClientPublisher    = 1
	ClientSubscriber   = 2
	UnAuthenticated    = 0
	Authenticated      = 1
)

func (c *SocketClient) startAuthenticate(store credentialStore) {

	log.Infof("waiting for client with id: %s to authenticate", c.clientId)

	msg, err := c.receive()

	if err != nil {
		c.close()
		return
	}

	authResult := c.authenticateWithCredentials(store, msg)

	if authResult {
		c.isAuthenticated = true
		c.sentAuthenticatedEvent()
	} else {
		c.close()
		return
	}

	typeResult := c.detectClientType(msg)

	if typeResult == 0 {
		c.close()
		return
	}

	c.clientType = typeResult

	c.startReceive()
}

func (c *SocketClient) authenticateWithCredentials(store credentialStore, request string) bool {

	cred := request[1:]

	if store.isValid(cred) {
		log.Infof("credential submitted: %s is valid", cred)
		return true
	} else {
		log.Warnf("credential submitted: %s is invalid", cred)
		return false
	}
}

func (c *SocketClient) detectClientType(request string) int {
	ctype, _ := strconv.Atoi(request[:1])

	switch ctype {
	case ClientPublisher:
		return ctype
	case ClientSubscriber:
		return ctype
	}

	log.Errorf("the client didn't provide a valid type, type: %s", request[:1])
	return 0
}

func (c *SocketClient) startReceive() {
	for {
		msg, err := c.receive()

		if err != nil {
			continue
		}

		c.onMessageChan <- &clientMessage{
			Payload:  msg,
			ClientId: c.clientId,
			Type:     clientMessageTypePublish,
		}

	}
}

func (c *SocketClient) receive() (string, error) {
	msg, err := read(c.connection, 1024)
	return string(msg), err
}

func (c *SocketClient) close() {
	_ = c.connection.Close()
	c.onMessageChan <- &clientMessage{ClientId: c.clientId, Type: clientMessageTypeDisconnect}
}

func (c *SocketClient) sentAuthenticatedEvent() {

}
