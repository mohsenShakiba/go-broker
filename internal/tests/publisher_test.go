package tests

import (
	"bufio"
	"fmt"
	"go-broker/internal/manager"
	"go-broker/internal/models"
	"net"
	"strconv"
	"sync"
	"testing"
	"time"
)

// this file will benchmark the performance of publishing messages to server
// it will use memory storage
func TestPublisherWithMemoryStorage(t *testing.T) {

	conf := manager.Config{
		StorageType: "M",
		Port:        8080,
	}

	_, err := manager.InitManager(conf)

	if err != nil {
		t.Fatalf("failed to create manager")
	}

	client, _ := net.Dial("tcp", "127.0.0.1:8080")
	clientReader := bufio.NewReader(client)

	var count int
	var lock sync.Mutex

	// ticker for rps
	go func() {
		ticker := time.NewTicker(time.Second)
		for {
			<-ticker.C
			fmt.Printf("RPS is %d\n", count)
			lock.Lock()
			count = 0
			lock.Unlock()
		}
	}()

	go func() {
		for {

			id := strconv.Itoa(count)

			ping := models.Message{
				Id:      id,
				Route:   "r1",
				Payload: []byte(id),
			}

			_ = ping.Write(client)

			_, _ = clientReader.ReadSlice('\n')

			ack := models.Ack{}
			_ = ack.FromReader(clientReader)

			if ack.Id != ping.Id {
				t.Fatal("the response of publish doesn't match")
			}

			lock.Lock()
			count += 1
			lock.Unlock()
		}
	}()

	time.Sleep(time.Second * 10)

}

// this file will benchmark the performance of publishing messages to server
// it will use file storage
func TestPublisherWithFileStorage(t *testing.T) {

	conf := manager.Config{
		FilePath:    "../files",
		StorageType: "F",
		Port:        8080,
	}

	_, err := manager.InitManager(conf)

	if err != nil {
		t.Fatalf("failed to create manager")
	}

	client, _ := net.Dial("tcp", "127.0.0.1:8080")
	clientReader := bufio.NewReader(client)

	var count int
	var lock sync.Mutex

	// ticker for rps
	go func() {
		ticker := time.NewTicker(time.Second)
		for {
			<-ticker.C
			fmt.Printf("RPS is %d\n", count)
			lock.Lock()
			count = 0
			lock.Unlock()
		}
	}()

	go func() {
		for {

			id := strconv.Itoa(count)

			ping := models.Message{
				Id:      id,
				Route:   "r1",
				Payload: []byte(id),
			}

			_ = ping.Write(client)

			_, _ = clientReader.ReadSlice('\n')

			ack := models.Ack{}
			_ = ack.FromReader(clientReader)

			if ack.Id != ping.Id {
				t.Fatal("the response of publish doesn't match")
			}

			lock.Lock()
			count += 1
			lock.Unlock()
		}
	}()

	time.Sleep(time.Second * 10)

}
