package socketserver

import (
	"fmt"
	"net"
	"strconv"
	"testing"
	"time"
)

// todo: change how receive works(use bytes instead of string)
// todo: use a buffer pool for reading from socket

func Test_Publisher_Valid(t *testing.T) {

	// valid auth
	validCred := "test1"

	// socket server config
	conf := SocketServerConfig{
		Credentials:    []string{validCred},
		ConnectionPort: 5555,
	}

	// socket server channel
	schan := make(chan string)

	// create socket server
	s := Init(conf, schan)

	// connect to server
	client, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", conf.ConnectionPort))

	if err != nil {
		t.Errorf("cannot connect to socket server, err: %s\n", err)
	}

	length := fmt.Sprintf("%04d", len(validCred)+5)
	authMsg := length + strconv.Itoa(1) + validCred

	t.Logf("authentication log for publisher is %s\n", authMsg)

	_, err = client.Write([]byte(authMsg))

	if err != nil {
		t.Errorf("error while authenitcating publisher, err: %s", err)
	}

	time.Sleep(time.Second * 1)

	socketClient := s.clients[0]

	if !socketClient.isAuthenticated {
		t.Errorf("the created client must be authenticated at this point")
	}
}

//func Test_Publisher_Invalid(t *testing.T) {
//
//}
//
//func Test_Subscriber_Valid(t *testing.T) {
//
//}
//
//func Test_Subscriber_Invalid(t *testing.T) {
//
//}
//
//func Test_Subscriber_Closed(t *testing.T) {
//
//}
//
//func TestRemoveClient(t *testing.T) {
//
//}
