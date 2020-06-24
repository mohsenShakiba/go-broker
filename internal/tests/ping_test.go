package tests

import (
	"bufio"
	"fmt"
	"go-broker/internal/manager"
	"go-broker/internal/models"
	"net"
	"sync"
	"testing"
	"time"
)

// this file will create a ping socket server
// it's main purpose is to message performance
// ping messages doesn't go through any channel
// and only measures the socket server performance
func TestPingServer(t *testing.T) {

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

			ping := models.Ping{
				Id: "IMAGINE_RANDOM_STRING",
			}

			_ = ping.Write(client)

			_, _ = clientReader.ReadSlice('\n')

			ack := models.Ack{}
			_ = ack.FromReader(clientReader)

			if ack.Id != ping.Id {
				t.Fatal("the response of ping doesn't match")
			}

			lock.Lock()
			count += 1
			lock.Unlock()

		}
	}()

	time.Sleep(time.Second * 10)

}
