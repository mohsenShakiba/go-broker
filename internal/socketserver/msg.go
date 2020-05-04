package socketserver

type clientMessage struct {
	ClientId string
	Type int
	Payload string
}

const (
	clientMessageTypeDisconnect = 1
	clientMessageTypePublish    = 2
)