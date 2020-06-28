package tests

import (
	"bufio"
	"fmt"
	"go-broker/internal/manager"
	"go-broker/internal/models"
	"go-broker/internal/storage"
	"net"
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestFull(t *testing.T) {

	conf := manager.Config{
		StorageConfig: storage.StorageConfig{
			Path:        "../files",
			Type:        "F",
			MaxFileSize: 1024 * 1024 * 1024,
		},
		Port: 8080,
	}

	_, err := manager.InitManager(conf)

	if err != nil {
		t.Fatalf("failed to create manager")
	}

	client, _ := net.Dial("tcp", "127.0.0.1:8080")
	clientReader := bufio.NewReader(client)

	var publishedMessageCount int
	var recievedMessageCount int
	var lock sync.Mutex

	// ticker for rps
	go func() {
		ticker := time.NewTicker(time.Second)
		for {
			<-ticker.C
			fmt.Printf("RPS, PUB: %d SUB: %d\n", publishedMessageCount, recievedMessageCount)
			lock.Lock()
			//publishedMessageCount = 0
			//recievedMessageCount = 0
			lock.Unlock()
		}
	}()

	// writing
	go func() {
		for {

			id := strconv.Itoa(publishedMessageCount)

			msg := models.Message{
				Id:      id,
				Route:   "r1",
				Payload: []byte(id),
			}

			_ = msg.Write(client)

			publishedMessageCount += 1
		}
	}()

	// reading
	go func() {
		for {
			_, _ = clientReader.ReadSlice('\n')

			ack := models.Ack{}
			_ = ack.FromReader(clientReader)
		}
	}()

	// creating subscriber
	subscriberClient, _ := net.Dial("tcp", "127.0.0.1:8080")
	subscriberClientReader := bufio.NewReader(subscriberClient)

	subMsg := &models.Register{
		Id:     "0",
		Dop:    1,
		Routes: []string{"r1"},
	}

	err = subMsg.Write(subscriberClient)

	if err != nil {
		t.Fatalf("failed to write message to server, err: %s", err)
	}

	_, _ = subscriberClientReader.ReadSlice('\n')

	ack := models.Ack{}
	_ = ack.FromReader(subscriberClientReader)

	go func() {
		for {
			msg := models.Message{}

			msgtype, _ := subscriberClientReader.ReadSlice('\n')

			if string(msgtype[:3]) != "PUB" {
				t.Fatalf("message type is invalid, error: %s", string(msgtype[:3]))
			}

			err = msg.FromReader(subscriberClientReader)

			if err != nil {
				t.Fatalf("failed to read payload msg, error: %s", err)
			}

			ackMsg := models.Ack{
				Id: msg.Id,
			}

			err = ackMsg.Write(subscriberClient)
			recievedMessageCount += 1

			if err != nil {
				t.Fatalf("failed to send ack messgae")
			}
		}
	}()

	time.Sleep(time.Second * 30)
}
