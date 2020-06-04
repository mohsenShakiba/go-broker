package tcp

import (
	"bufio"
	"github.com/google/uuid"
	"go-broker/internal/tcp/messages"
	"io"
	"net"
	"sync"
)

const (
	bufferSize = 1024
)

// subscriber is in charge of reading the data from the conn
// and sending data to connecting using the default protocol
type Client struct {
	ClientId string
	Reader   *bufio.Reader
	Writer   *bufio.Writer
	conn     io.ReadWriteCloser
	lock     sync.Mutex
}

// initSocketClient will create a new socket client
func initSocketClient(conn net.Conn) *Client {
	return &Client{
		ClientId: uuid.New().String(),
		Reader:   bufio.NewReader(conn),
		Writer:   bufio.NewWriter(conn),
		conn:     conn,
		lock:     sync.Mutex{},
	}
}

// Read will read an entire messages from the socket
func (c *Client) Read() (*messages.Message, bool) {
	c.lock.Lock()

	defer func() {
		c.lock.Unlock()
	}()

	return messages.ReadFromIO(c.Reader)
}

// Write will write the messages along with the prefix to the client connection
func (c *Client) Write(msg *messages.Message) {

	c.lock.Lock()

	defer func() {
		c.lock.Unlock()
	}()

	messages.WriteToIO(msg, c.Writer)
}

// Close will close the socket conn
func (c *Client) Close() error {

	c.lock.Lock()

	defer func() {
		c.lock.Unlock()
	}()

	err := c.conn.Close()

	return err
}
