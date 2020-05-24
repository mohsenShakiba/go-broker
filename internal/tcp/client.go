package tcp

import (
	"bufio"
	"github.com/google/uuid"
	"go-broker/internal/message"
	"io"
	"net"
	"sync"
)

const (
	bufferSize = 1024
)

// Client is in charge of reading the data from the conn
// and sending data to connecting using the default protocol
type Client struct {
	ClientId string
	reader   *bufio.Reader
	writer   *bufio.Writer
	conn     io.ReadWriteCloser
	lock     sync.Mutex
}

// initSocketClient will create a new socket client
func initSocketClient(conn net.Conn) *Client {
	return &Client{
		ClientId: uuid.New().String(),
		reader:   bufio.NewReaderSize(conn, bufferSize),
		writer:   bufio.NewWriterSize(conn, bufferSize),
		conn:     conn,
		lock:     sync.Mutex{},
	}
}

// Read will read an entire message from the socket
func (c *Client) Read() (*message.Message, bool) {
	c.lock.Lock()

	defer func() {
		c.lock.Unlock()
	}()

	return message.ReadFromIO(c.reader)
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

// Write will write the message along with the prefix to the client connection
func (c *Client) Write(msg *message.Message) {

	c.lock.Lock()

	defer func() {
		c.lock.Unlock()
	}()

	message.WriteToIO(msg, c.writer)
}
