package socketserver

import (
	"bytes"
	log "github.com/sirupsen/logrus"
	"strconv"
)

const (
	authenticateMessageType = "AUT"
	routedMessageType       = "PUB"
	subscribeMessageType    = "SUB"
	ackMessageType          = "ACK"
	nackMessageType         = "NCK"
)

//// this tcpMessage is sent to subscriber
//type contentMessage struct {
//	Id      string
//	Payload []byte
//}
//
//type receiveMessage struct {
//	Id      string
//	Success bool
//}
//
//// this tcpMessage is sent by publisher and subscribers for authentication
//type authenticateMessage struct {
//	Id       string
//	UserName string
//	Password string
//}
//
//// this tcpMessage is sent by the publisher to a route
//type routedMessage struct {
//	Id      string
//	Routes  []string
//	Payload []byte
//}
//
//// this tcpMessage is sent by the subscriber
//type subscribeMessage struct {
//	Id      string
//	Routes  []string
//	BufSize int
//}
//
//// this tcpMessage is sent by the subscriber to discard the tcpMessage as processed
//type ackMessage struct {
//	id string
//}
//
//// this tcpMessage is sent by the subscriber to requeue the tcpMessage
//type nackMessage struct {
//	id string
//}
//
