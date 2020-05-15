package socketserver

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
