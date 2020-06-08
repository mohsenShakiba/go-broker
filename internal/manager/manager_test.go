package manager

import (
	"bufio"
	log "github.com/sirupsen/logrus"
	"go-broker/internal/tcp/messages"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"testing"
	"time"
)

var numberOfSentMessages = 0
var numberOfReceivedMessages = 0
var counter = 0

func TestFull(t *testing.T) {

	dir, err := ioutil.TempDir("./", "temp")

	if err != nil {
		t.Fatal(err)
	}

	log.SetLevel(log.WarnLevel)

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
		msg := messages.NewMessage("PUB", strconv.Itoa(numberOfSentMessages))
		msg.WriteStr("routes", "r1")
		msg.WriteStr("payload", string(numberOfSentMessages))

		writer := bufio.NewWriterSize(publisherClient, 1)

		ok := messages.WriteToIO(msg, writer)

		writer.Flush()

		if !ok {
			t.Fatalf("could not write to server")
		}

		numberOfSentMessages += 1

		time.Sleep(time.Millisecond)
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
	subMsg.WriteStr("dop", "1")

	ok := messages.WriteToIO(subMsg, writer)

	if !ok {
		t.Fatalf("could not write to server")
	}

	ticker := time.NewTicker(time.Second)
	go func() {
		for {
			<-ticker.C
			t.Logf("RPS is %d", counter)
			counter = 0
		}
	}()

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

			ackMsg := messages.NewMessage("ACK", publishedMsg.MsgId)
			ok = messages.WriteToIO(ackMsg, writer)

			if !ok {
				t.Fatalf("failed to write ack message")
			}

			writer.Flush()
			counter += 1

		}
	}()

}
