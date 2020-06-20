package tcp

import "go-broker/internal/models"

type Context struct {
	Message interface{}
	Client  *Client
}

func (c *Context) SendError(id string, err string) {
	err := models.Err{
		Id: id,
	}
}
