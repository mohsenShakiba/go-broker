package tcp

import (
	"go-broker/internal/tcp/util"
	"net"
	"sync"
	"time"
)

type socketClient struct {
	clientId      string
	isClosed      bool
	connection    net.Conn
	onMessageChan chan<- clientMessage
	lock          sync.Mutex
}

type clientMessage struct {
	clientId string
	payload  []byte
}

func (c *socketClient) startReceive() {
	for {
		c.connection.SetReadDeadline(time.Now().Add(time.Second))
		c.read()
	}
}

func (c *socketClient) read() {
	c.lock.Lock()

	defer func() {
		c.lock.Unlock()
	}()

	if c.isClosed {
		return
	}

	msg, err := util.Read(c.connection, 1024)

	if err != nil {
		return
	}

	if msg == nil {
		return
	}

	c.onMessageChan <- clientMessage{
		clientId: c.clientId,
		payload:  msg,
	}
}

func (c *socketClient) Close() error {
	defer func() {
		c.isClosed = true
	}()
	err := c.connection.Close()
	return err
}

func (c *socketClient) Write(b []byte) (int, error) {
	c.lock.Lock()
	defer func() {
		c.lock.Unlock()
	}()
	return c.connection.Write(b)
}
