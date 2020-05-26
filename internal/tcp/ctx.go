package tcp

import "go-broker/internal/tcp/messages"

// Context contains information about tcp messages and it's context
type Context struct {
	Message *messages.Message
	Client  *Client
}
