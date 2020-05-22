package manager

import (
	"bytes"
	"fmt"
	log "github.com/sirupsen/logrus"
	"go-broker/internal/tcp/util"
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

	go initPublisher(t)
	go initSubscriber(t)

	time.Sleep(time.Second * 1000)
}

func initPublisher(t *testing.T) {
	publisherClient, err := net.Dial("tcp", "127.0.0.1:8085")

	if err != nil {
		t.Fatalf("could not establish client connection")
	}

	for {

		t.Logf("sending %d message", numberOfSentMessages)
		msg := fmt.Sprintf("type:PUB\nmsgId:%d\nroutes:r1\npayload:%d", numberOfSentMessages, numberOfSentMessages)
		l := fmt.Sprintf("%04d\n", len(msg)+5)
		_, err := publisherClient.Write([]byte(l))

		if err != nil {
			t.Fatalf("could not write to server, error: %s", err)
		}

		_, err = publisherClient.Write([]byte(msg))

		if err != nil {
			t.Fatalf("could not write to server, error: %s", err)
		}

		numberOfSentMessages += 1
		time.Sleep(time.Second)
	}
}

func initSubscriber(t *testing.T) {
	subscribeClient, err := net.Dial("tcp", "127.0.0.1:8085")

	if err != nil {
		t.Fatalf("could not establish client connection")
	}

	msg := fmt.Sprintf("type:SUB\nmsgId:%d\nroutes:r1", -1)
	l := fmt.Sprintf("%04d\n", len(msg)+5)
	_, err = subscribeClient.Write([]byte(l))

	if err != nil {
		t.Fatalf("could not write to server, error: %s", err)
	}

	_, err = subscribeClient.Write([]byte(msg))

	if err != nil {
		t.Fatalf("could not write to server, error: %s", err)
	}
	time.Sleep(time.Second * 2)

	go func() {
		for {

			msg, err := util.Read(subscribeClient, 1024)
			msgMap := make(map[string][]byte)

			if err != nil {
				t.Fatalf("could not read from socket client, err: %s", err)
			}

			newLineB := []byte("\n")
			colonB := []byte(":")

			// split by new line
			partsByNewLine := bytes.Split(msg, newLineB)

			for _, part := range partsByNewLine {
				partsByColon := bytes.Split(part, colonB)

				if len(partsByColon) != 2 {
					log.Warnf("bad payload data, discarding, message: %s", string(part))
				}

				msgMap[string(partsByColon[0])] = partsByColon[1]
			}

			msgType := msgMap["type"]
			msgId := msgMap["msgId"]

			if string(msgType) == "ACK" {
			} else {
				t.Logf("received message with id: %s", msgId)
			}

			time.Sleep(time.Second)
		}
	}()

}
