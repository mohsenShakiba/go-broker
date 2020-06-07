package manager

import (
	"bufio"
	"go-broker/internal/tcp/messages"
	"io/ioutil"
	"net"
	"os"
	"testing"
	"time"
)

var numberOfSentMessages = 0
var numberOfReceivedMessages = 0

func TestFull(t *testing.T) {

	dir, err := ioutil.TempDir("./", "temp")

	if err != nil {
		t.Fatal(err)
	}

	defer os.RemoveAll(dir)

	_, err = InitManager(dir)

	if err != nil {
		t.Fatalf("error while creating manager, error: %s", err)
	}

	go initSubscriber(t)

	time.Sleep(time.Second)

	go initPublisher(t)

	time.Sleep(time.Second * 100000)
}

func initPublisher(t *testing.T) {
	publisherClient, err := net.Dial("tcp", "127.0.0.1:8085")

	if err != nil {
		t.Fatalf("could not establish client connection")
	}

	for {
		msg := messages.NewMessage("PUB", string(numberOfSentMessages))
		msg.WriteStr("routes", "route")
		msg.WriteStr("payload", string(numberOfSentMessages))

		writer := bufio.NewWriterSize(publisherClient, 1)

		ok := messages.WriteToIO(msg, writer)

		if !ok {
			t.Fatalf("could not write to server")
		}

		numberOfSentMessages += 1

		t.Logf("sent message with id: %d", numberOfSentMessages)

		time.Sleep(time.Millisecond * 1000)
	}
}

func initSubscriber(t *testing.T) {

	subscribeClient, err := net.Dial("tcp", "127.0.0.1:8085")
	writer := bufio.NewWriterSize(subscribeClient, 1)
	reader := bufio.NewReaderSize(subscribeClient, 1)

	if err != nil {
		t.Fatalf("could not establish client connection")
	}

	// send subscription message

	subMsg := messages.NewMessage("SUB", "-")
	subMsg.WriteStr("routes", "r1")
	subMsg.WriteStr("dop", "10")

	ok := messages.WriteToIO(subMsg, writer)

	if !ok {
		t.Fatalf("could not write to server")
	}

	go func() {

		// read the subscription result
		subResMsg, ok := messages.ReadFromIO(reader)

		if !ok {
			t.Fatalf("could not read from server")
		}

		if subResMsg.Type != "ACK" {
			t.Fatalf("the message type must be ack")
		}

		for {
			publishedMsg, ok := messages.ReadFromIO(reader)

			if !ok {
				t.Fatalf("could not read from server")
			}

			t.Logf("received message with id: %s", publishedMsg.MsgId)
		}
	}()

}
