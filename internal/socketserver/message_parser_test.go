package socketserver

import (
	"encoding/binary"
	"strings"
	"testing"
)

func TestAuthMessage(t *testing.T) {
	inMsg := authenticateMessage{
		id:       "TEST",
		userName: "USER",
		password: "PASS",
	}

	b := format([]byte(authenticateMessageType), []byte(inMsg.id), []byte(inMsg.userName), []byte(inMsg.password))

	outMsg := parseMessage(b)

	switch msg := outMsg.(type) {
	case authenticateMessage:
		if msg.userName != inMsg.userName || msg.password != inMsg.password || msg.id != inMsg.id {
			t.Fatalf("the output properties deosn't match that of input")
		}
	default:
		t.Fatalf("the output message type was not valid")
	}

}

func TestRoutedMessage(t *testing.T) {
	inMsg := routedMessage{
		id:      "TEST",
		payload: []byte("USER"),
		routes:  []string{"ROUTE1", "ROUTE2"},
	}

	b := format([]byte(routedMessageType), []byte(inMsg.id), []byte(strings.Join(inMsg.routes, ",")), inMsg.payload)

	outMsg := parseMessage(b)

	switch msg := outMsg.(type) {
	case routedMessage:
		if string(msg.payload) != string(inMsg.payload) || strings.Join(msg.routes, ",") != strings.Join(inMsg.routes, ",") || msg.id != inMsg.id {
			t.Fatalf("the output properties deosn't match that of input")
		}
	default:
		t.Fatalf("the output message type was not valid")
	}

}

// test subscribe message
func TestSubscribeMessage(t *testing.T) {
	inMsg := subscribeMessage{
		id:      "TEST",
		routes:  []string{"ROUTE1", "ROUTE2"},
		bufSize: 10,
	}

	b := format([]byte(routedMessageType), []byte(inMsg.id), []byte(strings.Join(inMsg.routes, ",")))

	outMsg := parseMessage(b)

	switch msg := outMsg.(type) {
	case routedMessage:
		if string(msg.payload) != string(inMsg.payload) || strings.Join(msg.routes, ",") != strings.Join(inMsg.routes, ",") || msg.id != inMsg.id {
			t.Fatalf("the output properties deosn't match that of input")
		}
	default:
		t.Fatalf("the output message type was not valid")
	}

}

// test ack message

// test nack message
