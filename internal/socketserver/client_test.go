package socketserver

import (
	"fmt"
	"testing"
)

func TestClientAuthenticateValid(t *testing.T) {
	validCred := "Cred"
	client := socketClient{}
	store := credentialStore{
		config: SocketServerConfig{
			Credentials: []string{validCred},
		},
	}

	validCredWithType := fmt.Sprintf("%d%s", clientPublisher, validCred)

	authSuccess := client.authenticate(store, validCredWithType)

	if !authSuccess {
		t.Fatal("client authentication must succeed")
	}

}

func TestClientAuthenticateInvalid(t *testing.T) {
	validCred := "Cred"
	invalidCred := "DifferentCred"
	client := socketClient{}
	store := credentialStore{
		config: SocketServerConfig{
			Credentials: []string{validCred},
		},
	}

	invalidCredWithType := fmt.Sprintf("%d%s", clientPublisher, invalidCred)

	authSuccess := client.authenticate(store, invalidCredWithType)

	if authSuccess {
		t.Fatal("client authentication must fail")
	}

}

func TestClientTypeValid(t *testing.T) {
	validType := clientPublisher
	client := socketClient{}

	input := fmt.Sprintf("%d", validType)

	ctype, _ := client.parseClientType(input)

	if ctype != validType {
		t.Fatalf("invalid output while detecting client type, input:%d, output:%d", validType, ctype)
	}
}

func TestClientTypeInvalid(t *testing.T) {
	validType := 5
	client := socketClient{}

	input := fmt.Sprintf("%d", validType)

	_, ok := client.parseClientType(input)

	if ok {
		t.Fatalf("the client type must detect and error, input:%d", validType)
	}
}
