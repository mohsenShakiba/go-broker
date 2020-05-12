package socketserver

import (
	"go-broker/internal/socketserver/util"
	"io"
)

type socketClient struct {
	clientId        string
	isAuthenticated bool
	clientType      int
	isClosed        bool
	connection      io.ReadWriteCloser
	onMessageChan   chan<- clientMessage
}

type clientMessage struct {
	clientId string
	payload  []byte
}

const (
	clientUndetermined = 0
	clientPublisher    = 1
	clientSubscriber   = 2
)

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

		msg, err := util.Read(c.connection, 1024)

		if err != nil {
			continue
		}

		c.onMessageChan <- clientMessage{
			clientId: c.clientId,
			payload:  msg,
		}

	}
}

func (c *socketClient) close() {
	_ = c.connection.Close()
	c.isClosed = true
}

func (c *socketClient) send(payload []byte) error {
	_, err := c.connection.Write(payload)
	return err
}
