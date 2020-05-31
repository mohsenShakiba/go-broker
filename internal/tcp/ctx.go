package tcp

import "go-broker/internal/tcp/messages"

// Context contains information about tcp messages and it's context
type Context struct {
	Message *messages.Message
	Client  *Client
}

func (c *Context) SendAck() {
	m := messages.NewMessage("ACK", c.Message.MsgId)
	c.Client.Write(m)
}

func (c *Context) SendErr(code string) {
	m := messages.NewMessage("ERR", c.Message.MsgId)
	m.WriteStr("code", code)
	c.Client.Write(m)
}
