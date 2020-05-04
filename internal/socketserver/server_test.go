package socketserver

import (
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"testing"
	"time"
)

// todo: change how receive works(use bytes instead of string)
// todo: use a buffer pool for reading from socket

func Test_Read_Small(t *testing.T) {
	body := "THIS IS A TEST"
	msg := toByteWithLengthPrefix(body)
	r := strings.NewReader(string(msg))
	result, err := read(r, 100)

	if err != nil {
		t.Errorf("error in reading, input: %s, error: %s", body, err)
	}

	if string(result) != body {
		t.Errorf("input and output mismatch while reading, input %s, output: %s", result, body)
	}
}

func Test_Read_Large(t *testing.T) {
	body := createRandomString(20)
	msg := toByteWithLengthPrefix(body)
	r := strings.NewReader(string(msg))
	result, err := read(r, 10)

	if err != nil {
		t.Errorf("error in reading, input: %s, error: %s", body, err)
	}

	if string(result) != body {
		t.Errorf("input and output mismatch while reading, input %s, output: %s", result, body)
	}
}

func Test_ClientAuthenticate_Valid(t *testing.T) {
	validCred := "Cred"
	client := SocketClient{}
	store := credentialStore{
		config: SocketServerConfig{
			Credentials: []string {validCred},
		},
	}

	validCredWithType := fmt.Sprintf("%d%s", ClientPublisher, validCred)

	authSuccess := client.authenticateWithCredentials(store, validCredWithType)

	if !authSuccess {
		t.Error("client authentication must succeed")
	}

}

func Test_ClientAuthenticate_Invalid(t *testing.T) {
	validCred := "Cred"
	invalidCred := "DifferentCred"
	client := SocketClient{}
	store := credentialStore{
		config: SocketServerConfig{
			Credentials: []string {validCred},
		},
	}

	invalidCredWithType := fmt.Sprintf("%d%s", ClientPublisher, invalidCred)

	authSuccess := client.authenticateWithCredentials(store, invalidCredWithType)

	if authSuccess {
		t.Error("client authentication must fail")
	}

}

func Test_ClientType_Valid(t *testing.T) {
	validType := ClientPublisher
	client := SocketClient{}

	input := fmt.Sprintf("%d", validType)

	output := client.detectClientType(input)

	if output != validType {
		t.Errorf("invalid output while detecting client type, input:%d, output:%d", validType, output)
	}
}

func Test_ClientType_Invalid(t *testing.T) {
	validType := 5
	client := SocketClient{}

	input := fmt.Sprintf("%d", validType)

	output := client.detectClientType(input)

	if output != 0 {
		t.Errorf("the client type must detect and error, input:%d, output:%d", validType, output)
	}
}


func Test_Publisher_Valid(t *testing.T) {

	// valid auth
	validCred := "test1"

	// socket server config
	conf := SocketServerConfig{
		Credentials: []string {validCred},
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

	length := fmt.Sprintf("%04d", len(validCred) + 5)
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

func Test_Publisher_Invalid(t *testing.T) {

}

func Test_Subscriber_Valid(t *testing.T) {

}

func Test_Subscriber_Invalid(t *testing.T) {

}

func Test_Subscriber_Closed(t *testing.T) {

}




func createRandomString(n int) string {
	letterRunes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
