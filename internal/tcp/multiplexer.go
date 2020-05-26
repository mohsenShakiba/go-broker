package tcp

import (
	log "github.com/sirupsen/logrus"
	"go-broker/internal/tcp/messages"
)

// Multiplexer is in charge of storing the mapping between handlers and contexts
type Multiplexer struct {
	handlers map[string]func(ctx *Context)
}

// creates a new Multiplexer
func newMultiplexer() *Multiplexer {
	return &Multiplexer{
		handlers: make(map[string]func(ctx *Context)),
	}
}

// RegisterHandler will register a handler for the given type
// handlers with matching type will receive the messages
func (m *Multiplexer) RegisterHandler(t string, h func(msg *Context)) {
	m.handlers[t] = h
}

// process will forward the incoming messages to the appropriate handler
func (m *Multiplexer) process(client *Client, msg *messages.Message) {

	handler, ok := m.handlers[msg.Type]

	if !ok {
		log.Errorf("no handler was found for %m", msg.Type)
	}

	go handler(&Context{
		Message: msg,
		Client:  client,
	})

}
