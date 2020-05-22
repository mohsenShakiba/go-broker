package storage

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestWrite(t *testing.T) {

	dir, err := ioutil.TempDir("./", "temp")

	if err != nil {
		t.Fatal(err)
	}

	defer os.RemoveAll(dir)

	s, err := Init(dir)

	if err != nil {
		t.Fatalf("could not create a new storage instance, err: %s", err)
	}

	msg := Message{
		MsgId:   "test",
		Routes:  []string{"r1", "r2"},
		Payload: []byte("THIS IS A TEST"),
	}

	s.Add(&msg)

	s, err = Init(dir)

	if err != nil {
		t.Fatalf("could not recreate a new storage %s", err)
	}

	response := s.read("test")

	if string(response.Payload) != string(msg.Payload) {
		t.Fatalf("the Payload of messages don't match")
	}

}
