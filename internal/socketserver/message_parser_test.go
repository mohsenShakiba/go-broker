package socketserver

import (
	"bytes"
	"go-broker/internal/socketserver/serializer"
	"strings"
	"testing"
)

func TestAuthMessage(t *testing.T) {
	inMsg := authenticateMessage{
		Id:       "TEST",
		UserName: "USER",
		Password: "PASS",
	}

	s := serializer.NewJsonSerializer()

	b, err := s.Serialize(&inMsg)

	if err != nil {
		t.Fatalf("serialization failed with error %s", err)
	}

	buf := new(bytes.Buffer)
	buf.Write([]byte(authenticateMessageType))
	buf.Write(b)

	outMsg, err := parseMessage(buf.Bytes())

	if err != nil {
		t.Fatalf("failed to parse message, error: %s", err)
	}

	switch msg := outMsg.(type) {
	case *authenticateMessage:
		if msg.UserName != inMsg.UserName || msg.Password != inMsg.Password || msg.Id != inMsg.Id {
			t.Fatalf("the output properties deosn't match that of input")
		}
	default:
		t.Fatalf("the output message type was not valid")
	}

}

func TestRoutedMessage(t *testing.T) {
	inMsg := routedMessage{
		Id:      "TEST",
		Payload: []byte("USER"),
		Routes:  []string{"ROUTE1", "ROUTE2"},
	}

	s := serializer.NewJsonSerializer()

	b, err := s.Serialize(&inMsg)

	if err != nil {
		t.Fatalf("serialization failed with error %s", err)
	}

	buf := new(bytes.Buffer)
	buf.Write([]byte(routedMessageType))
	buf.Write(b)

	outMsg, err := parseMessage(buf.Bytes())

	if err != nil {
		t.Fatalf("failed to parse message, error: %s", err)
	}

	switch msg := outMsg.(type) {
	case *routedMessage:
		if string(msg.Payload) != string(inMsg.Payload) || strings.Join(msg.Routes, ",") != strings.Join(inMsg.Routes, ",") || msg.Id != inMsg.Id {
			t.Fatalf("the output properties deosn't match that of input")
		}
	default:
		t.Fatalf("the output message type was not valid")
	}

}

// test subscribe message
func TestSubscribeMessage(t *testing.T) {
	inMsg := subscribeMessage{
		Id:      "TEST",
		Routes:  []string{"ROUTE1", "ROUTE2"},
		BufSize: 10,
	}

	s := serializer.NewJsonSerializer()

	b, err := s.Serialize(inMsg)

	if err != nil {
		t.Fatalf("serialization failed with error %s", err)
	}

	buf := new(bytes.Buffer)
	buf.Write([]byte(subscribeMessageType))
	buf.Write(b)

	outMsg, err := parseMessage(buf.Bytes())

	if err != nil {
		t.Fatalf("failed to parse message, error: %s", err)
	}

	switch msg := outMsg.(type) {
	case *subscribeMessage:
		if msg.BufSize != inMsg.BufSize || strings.Join(msg.Routes, ",") != strings.Join(inMsg.Routes, ",") || msg.Id != inMsg.Id {
			t.Fatalf("the output properties deosn't match that of input")
		}
	default:
		t.Fatalf("the output message type was not valid")
	}

}

// test ack message

// test nack message
