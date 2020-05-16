package socketserver

type MessageHandler interface {
	HandleMessage(m *MessageContext)
}
