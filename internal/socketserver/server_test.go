package socketserver

import (
	"fmt"
	serializer "go-broker/internal/socketserver/serializer"
	"io"
	"net"
	"testing"
	"time"
)

func Test_Publisher_Valid(t *testing.T) {

	// setup
	validUserName := "USER"
	validPassword := "PASS"

	inValidUserName := "IUSER"
	invalidPassword := "IPASS"

	validCred := fmt.Sprintf("%s:%s", validUserName, validPassword)

	conf := SocketServerConfig{
		Credentials:    []string{validCred},
		ConnectionPort: 5555,
	}

	publishMessageChan := make(chan ServerEvents)

	// create new server
	_ = Init(conf, publishMessageChan)

	// create new publisher client
	publisherOne, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", conf.ConnectionPort))

	if err != nil {
		t.Fatalf("cannot connect to socket server for publisher one, err: %s\n", err)
	}

	// create new subscriber client
	subscriberOne, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", conf.ConnectionPort))

	if err != nil {
		t.Fatalf("cannot connect to socket server for subscriber, err: %s\n", err)
	}

	// test invalid authentication for publisher
	invalidAuthMsg := createAuthMessage(t, inValidUserName, invalidPassword)
	_, err = publisherOne.Write(invalidAuthMsg)

	if err != nil {
		t.Fatalf("failed to write msg to socker")
	}

	// must receive invalid auth event
	time.Sleep(time.Second)
	publisherOneRecEv := readAuthenticateMessage(t, publisherOne)

	if publisherOneRecEv.Success {
		t.Fatalf("the received authentication tcpMessage must indicate false result %t", publisherOneRecEv.Success)
	}

	// test invalid authentication for subscriber
	_, err = subscriberOne.Write(invalidAuthMsg)

	if err != nil {
		t.Fatalf("failed to write msg to socker")
	}

	// must receive invalid auth event
	time.Sleep(time.Second)
	subscriberOneRecEv := readAuthenticateMessage(t, subscriberOne)

	if subscriberOneRecEv.Success {
		t.Fatalf("the received authentication tcpMessage must indicate false result %t", subscriberOneRecEv.Success)
	}

	// send tcpMessage from unauthenticated publisher make sure tcpMessage is not received
	publishMessage(t, "1", "TEST1", "BODY")

	time.Sleep(time.Second)

	select {
	case <-publishMessageChan:
		t.Fatalf("we should not receive the publish tcpMessage as this point")
	default:
		break
	}

	// test valid authentication tcpMessage for publisher
	validAuthMsg := createAuthMessage(t, validUserName, validPassword)
	_, err = publisherOne.Write(validAuthMsg)

	time.Sleep(time.Second)

	publisherOneRecEv = readAuthenticateMessage(t, publisherOne)

	if !publisherOneRecEv.Success {
		t.Fatalf("the received authentication tcpMessage must indicate true result %t", publisherOneRecEv.Success)
	}

	// test valid authentication tcpMessage for subscriber
	_, err = subscriberOne.Write(validAuthMsg)

	time.Sleep(time.Second)

	subscriberOneRecEv = readAuthenticateMessage(t, subscriberOne)

	if !subscriberOneRecEv.Success {
		t.Fatalf("the received authentication tcpMessage must indicate true result %t", publisherOneRecEv.Success)
	}

	// send subscription tcpMessage
	subscriptionMessage := createSubscriptionMessage(t, "NOT_EXISTING_ROUTE")
	_, err = subscriberOne.Write(subscriptionMessage)

	if err != nil {
		t.Fatalf("failed to send subscription tcpMessage, error: %s", err)
	}

	// must receive subscription event

	// test publish a tcpMessage to invalid route

	// test receive tcpMessage from invalid route

	// test publish a tcpMessage to valid route

	// test receive tcpMessage from valid route

	// test nack tcpMessage

	// create new subscriber client

	// tcpMessage must be received in the new client

	// test ack tcpMessage

	// create new subscriber

	// the tcpMessage must not be delivered to the new client

}

func createAuthMessage(t *testing.T, userName string, password string) []byte {
	msg := authenticateMessage{
		Id:       "-",
		UserName: userName,
		Password: password,
	}

	jsonSerializer := serializer.NewJsonSerializer()

	serializedMsg, err := jsonSerializer.Serialize(msg)

	if err != nil {
		t.Fatalf("serialization of authentication tcpMessage failed, error: %s", err)
	}

	return formatter("AUT", serializedMsg)
}

func readAuthenticateMessage(t *testing.T, reader io.Reader) *receiveMessage {
	jsonSerializer := serializer.NewJsonSerializer()
	b := make([]byte, 1024)
	_, err := reader.Read(b)

	if err != nil {
		t.Fatalf("failed to read data from binary %s", err)
	}

	msg := &receiveMessage{}
	jsonSerializer.Deserialize(b, msg)
	return msg
}

func createSubscriptionMessage(t *testing.T, route string) []byte {
	msg := subscribeMessage{
		Id:      "-",
		Routes:  []string{route},
		BufSize: 1,
	}

	jsonSerializer := serializer.NewJsonSerializer()

	serializedMsg, err := jsonSerializer.Serialize(msg)

	if err != nil {
		t.Fatalf("serialization of subscription tcpMessage failed, error: %s", err)
	}

	return formatter("SUB", serializedMsg)
}

func publishMessage(t *testing.T, id string, route string, payload string) []byte {
	msg := routedMessage{
		Id:      id,
		Routes:  []string{route},
		Payload: []byte(payload),
	}

	jsonSerializer := serializer.NewJsonSerializer()

	serializedMsg, err := jsonSerializer.Serialize(msg)

	if err != nil {
		t.Fatalf("serialization of authentication tcpMessage failed, error: %s", err)
	}

	return formatter("PUB", serializedMsg)
}

func formatter(t string, p []byte) []byte {
	prefix := []byte(fmt.Sprintf("%04d%s", len(p)+7, t))
	return append(prefix, p...)
}
