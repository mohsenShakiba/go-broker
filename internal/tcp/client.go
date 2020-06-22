package tcp

import (
	"bufio"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"go-broker/internal/models"
	"io"
	"net"
	"sync"
)

// subscriber is in charge of reading the data from the conn
// and sending data to connecting using the default protocol
type Client struct {
	ClientId string
	Reader   *bufio.Reader
	Writer   *bufio.Writer
	conn     io.ReadWriteCloser
	wLock    sync.Mutex
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

// BeginWrite will only allow one write at any time
// this is required to prevent multiple write
func (c *Client) BeginWrite(fn func(w io.Writer)) {
	//fn(c.conn)
	c.wLock.Lock()
	fn(c.Writer)
	c.Writer.Flush()
	//c.conn.Write(buf.Bytes())
	c.wLock.Unlock()

}

func (c *Client) SendError(id string, msg string) {
	e := models.Err{
		Id:  id,
		Err: msg,
	}

	c.BeginWrite(func(w io.Writer) {
		err := e.Write(w)
		if err != nil {
			log.Errorf("failed to write error to client, err: %s", err)
		}
	})

}

func (c *Client) SendAck(id string) {
	ack := models.Ack{
		Id: id,
	}

	c.BeginWrite(func(w io.Writer) {
		err := ack.Write(w)
		if err != nil {
			log.Errorf("failed to write ack to client, err: %s", err)
		}
	})
}

func (c *Client) SendNack(id string) {
	nack := models.Nack{
		Id: id,
	}

	c.BeginWrite(func(w io.Writer) {
		err := nack.Write(w)
		if err != nil {
			log.Errorf("failed to write nack to client, err: %s", err)
		}
	})

}

// Close will close the socket conn
func (c *Client) Close() error {

	err := c.conn.Close()

	return err
}
