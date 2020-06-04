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

	//go initPublisher(t)
	go initSubscriber(t)

	time.Sleep(time.Second * 100000)
}

func initPublisher(t *testing.T) {
	publisherClient, err := net.Dial("tcp", "127.0.0.1:8085")

	if err != nil {
		t.Fatalf("could not establish client connection")
	}

	for {
		msg := messages.NewMessage("PUB", string(numberOfSentMessages))
		msg.WriteStr("routes", "r1")
		msg.WriteStr("payload", string(numberOfSentMessages))

		writer := bufio.NewWriter(publisherClient)

		ok := messages.WriteToIO(msg, writer)

		if !ok {
			t.Fatalf("could not write to server")
		}

		numberOfSentMessages += 1

		t.Logf("sent message with id: %d", numberOfSentMessages)

		time.Sleep(time.Second)
	}
}

func initSubscriber(t *testing.T) {

	subscribeClient, err := net.Dial("tcp", "127.0.0.1:8085")
	writer := bufio.NewWriter(subscribeClient)

	if err != nil {
		t.Fatalf("could not establish client connection")
	}

	// send subscription message

	subMsg := messages.NewMessage("SUB", "-")
	subMsg.WriteStr("routes", "r1")

	written, err := writer.Write([]byte("fdsfsd\n"))

	t.Logf("written %d", written)

	ok := messages.WriteToIO(subMsg, writer)

	if !ok {
		t.Fatalf("could not write to server")
	}

	//// read the subscription result
	//subResMsg, ok := messages.ReadFromIO(reader)
	//
	//if !ok {
	//	t.Fatalf("could not read from server")
	//}
	//
	//if subResMsg.Type != "ACK" {
	//	t.Fatalf("the message type must be ack")
	//}
	//
	//go func() {
	//	for {
	//		publishedMsg, ok := messages.ReadFromIO(reader)
	//
	//		if !ok {
	//			t.Fatalf("could not read from server")
	//		}
	//
	//		t.Logf("received message with id: %s", publishedMsg.MsgId)
	//	}
	//}()

}
