package tcp

import (
	"bufio"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"go-broker/internal/models"
	"io"
	"net"
)

// subscriber is in charge of reading the data from the conn
// and sending data to connecting using the default protocol
type Client struct {
	ClientId string
	Reader   *bufio.Reader
	Writer   *bufio.Writer
	conn     io.ReadWriteCloser
}

// initSocketClient will create a new socket client
func initSocketClient(conn net.Conn) *Client {
	return &Client{
		ClientId: uuid.New().String(),
		Reader:   bufio.NewReader(conn),
		Writer:   bufio.NewWriter(conn),
		conn:     conn,
	}
}

// Read will read an entire messages from the socket
func (c *Client) Read() (interface{}, error) {
	return models.Parse(c.Reader)
}

func (c *Client) Write(b []byte) (int, error) {
	return c.conn.Write(b)
}

func (c *Client) SendError(id string, msg string) {
	e := models.Err{
		Id:  id,
		Err: msg,
	}

	err := e.Write(c)

	log.Errorf("error: %s", msg)

	if err != nil {
		log.Errorf("failed to write error to client, err: %s", err)
	}

}

func (c *Client) SendAck(id string) {
	ack := models.Ack{
		Id: id,
	}

	err := ack.Write(c)

	if err != nil {
		log.Errorf("failed to write ack to client, err: %s", err)
	}
}

func (c *Client) SendNack(id string) {
	ack := models.Nack{
		Id: id,
	}

	err := ack.Write(c)

	if err != nil {
		log.Errorf("failed to write nack to client, err: %s", err)
	}

}

// Close will close the socket conn
func (c *Client) Close() error {

	err := c.conn.Close()

	return err
}
