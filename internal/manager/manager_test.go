package manager

import (
	"bufio"
	"go-broker/internal/models"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"
)

var numberOfSentMessages = 0
var numberOfReceivedMessages = 0
var counter = 0
var pub_counter = 0

func TestFull(t *testing.T) {

	dir, err := ioutil.TempDir("./", "temp")

	if err != nil {
		t.Fatal(err)
	}

	//log.SetLevel(log.WarnLevel)

	defer os.RemoveAll(dir)

	conf := Config{
		FilePath:    dir,
		StorageType: "F",
		Port:        8080,
	}

	_, err = InitManager(conf)

	if err != nil {
		t.Fatalf("error while creating manager, error: %s", err)
	}

	go initSubscriber(t)

	time.Sleep(time.Second)

	go initPublisher(t)

	time.Sleep(time.Second * 100000)
}

func initPublisher(t *testing.T) {
	publisherClient, err := net.Dial("tcp", "127.0.0.1:8080")

	if err != nil {
		t.Fatalf("could not establish client connection")
	}

	ticker := time.NewTicker(time.Second)
	go func() {
		for {
			<-ticker.C
			//t.Logf("RPS PUB is %d", pub_counter)
			//pub_counter = 0
		}
	}()

	for {
		numberOfSentMessages += 1

		msg := models.Message{
			Id:      strconv.Itoa(numberOfSentMessages),
			Route:   "t1",
			Payload: []byte(strconv.Itoa(numberOfSentMessages)),
		}

		err := msg.Write(publisherClient)

		if err != nil {
			t.Fatalf("could not write to server")
		}

		pub_counter += 1

		//time.Sleep(time.Millisecond * 1)
	}
}

func initSubscriber(t *testing.T) {

	subscribeClient, err := net.Dial("tcp", "127.0.0.1:8080")
	reader := bufio.NewReader(subscribeClient)

	if err != nil {
		t.Fatalf("could not establish client connection")
	}

	// send subscription message

	subMsg := &models.Register{
		Id:     "0",
		Dop:    9999,
		Routes: []string{"t1"},
	}

	err = subMsg.Write(subscribeClient)

	if err != nil {
		t.Fatalf("failed to write message to server, err: %s", err)
	}

	ticker := time.NewTicker(time.Second)
	go func() {
		for {
			<-ticker.C
			t.Logf("RPS SUB is %d", counter)
			counter = 0
		}
	}()

	go func() {

		// read the subscription ack
		subAck := models.Ack{}

		msgtype, _ := reader.ReadSlice('\n')

		if string(msgtype[:3]) != "ACK" {
			t.Fatalf("message type is invalid, error: %s", string(msgtype[:3]))
		}

		err := subAck.FromReader(reader)

		if err != nil {
			t.Fatalf("failed to read subscription ack, error: %s", err)
		}

		if subAck.Id != "0" {
			t.Fatalf("the message id doesn't match")
		}

		var l sync.Mutex

		for {
			msg := models.Message{}

			msgtype, _ := reader.ReadSlice('\n')

			if string(msgtype[:3]) != "PUB" {
				t.Fatalf("message type is invalid, error: %s", string(msgtype[:3]))
			}

			err = msg.FromReader(reader)

			if err != nil {
				t.Fatalf("failed to read payload msg, error: %s", err)
			}

			go func() {
				ackMsg := models.Ack{
					Id: msg.Id,
				}

				l.Lock()
				err = ackMsg.Write(subscribeClient)
				l.Unlock()

				if err != nil {
					t.Fatalf("failed to send ack messgae")
				}

				counter += 1
			}()

		}
	}()

}
