package socketserver

const (
	authenticateMessageType = "AUT"
	routedMessageType       = "PUB"
	subscribeMessageType    = "SUB"
	ackMessageType          = "ACK"
	nackMessageType         = "NCK"
)

// this message is sent to subscriber
type contentMessage struct {
	Id      string
	Payload []byte
}

type receiveMessage struct {
	Id      string
	Success bool
}

// this message is sent by publisher and subscribers for authentication
type authenticateMessage struct {
	Id       string
	UserName string
	Password string
}

// this message is sent by the publisher to a route
type routedMessage struct {
	Id      string
	Routes  []string
	Payload []byte
}

// this message is sent by the subscriber
type subscribeMessage struct {
	Id      string
	Routes  []string
	BufSize int
}

// this message is sent by the subscriber to discard the message as processed
type ackMessage struct {
	id string
}

// this message is sent by the subscriber to requeue the message
type nackMessage struct {
	id string
}

type ServerEvents interface{}

type SubscriberAuthenticatedEvent struct {
	ClientId string
	Routes   []string
	BufSize  int
}

type PublishMessageEvent struct {
	ClientId string
	MsgId    string
	Routes   []string
	Payload  []byte
}

type AckMessageEvent struct {
	ClientId string
	MsgId    string
}

type NackMessageEvent struct {
	ClientId string
	MsgId    string
}
