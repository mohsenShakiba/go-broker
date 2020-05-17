package tcp

import (
	"go-broker/internal/tcp/util"
	"io"
	"sync"
)

type socketClient struct {
	clientId      string
	isClosed      bool
	connection    io.ReadWriteCloser
	onMessageChan chan<- clientMessage
	lock          sync.Mutex
}

type clientMessage struct {
	clientId string
	payload  []byte
}

func (c *socketClient) startReceive() {
	for {
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

	c.onMessageChan <- clientMessage{
		clientId: c.clientId,
		payload:  msg,
	}

	c.lock.Unlock()
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
	return c.Write(b)
}
