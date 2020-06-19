package tcp

import (
	"bufio"
	"github.com/google/uuid"
	"go-broker/internal/models"
	"io"
	"net"
)

// subscriber is in charge of reading the data from the Conn
// and sending data to connecting using the default protocol
type Client struct {
	ClientId string
	Reader   *bufio.Reader
	Writer   *bufio.Writer
	Conn     io.ReadWriteCloser
}

// initSocketClient will create a new socket client
func initSocketClient(conn net.Conn) *Client {
	return &Client{
		ClientId: uuid.New().String(),
		Reader:   bufio.NewReader(conn),
		Writer:   bufio.NewWriter(conn),
		Conn:     conn,
	}
}

// Read will read an entire messages from the socket
func (c *Client) Read() (interface{}, error) {
	return models.Parse(c.Reader)
}

func (c *Client) Write(b []byte) (int, error) {
	return c.Conn.Write(b)
}

// Close will close the socket Conn
func (c *Client) Close() error {

	err := c.Conn.Close()

	return err
}
