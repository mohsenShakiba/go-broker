package socketserver

import (
	"bytes"
	"fmt"
	"math/rand"
	"strings"
	"testing"
)

func TestReadSmall(t *testing.T) {
	body := []byte("THIS IS A TEST")
	l := messageLengthPrefixFormatter(len(body))
	payload := fmt.Sprintf("%s%s", l, body)
	r := strings.NewReader(payload)
	result, err := read(r, 100)

	if err != nil {
		t.Fatalf("error in reading, input: %s, error: %s", body, err)
	}

	if string(result) != string(body) {
		t.Fatalf("input and output mismatch while reading, input %s, output: %s", result, body)
	}
}

func TestReadLarge(t *testing.T) {
	body := createRandomString(20)
	l := messageLengthPrefixFormatter(len(body))
	payload := fmt.Sprintf("%s%s", l, body)
	r := strings.NewReader(payload)
	result, err := read(r, 10)

	if err != nil {
		t.Fatalf("error in reading, input: %s, error: %s", body, err)
	}

	if string(result) != body {
		t.Fatalf("input and output mismatch while reading, input %s, output: %s", result, body)
	}
}

func TestReadInvalid(t *testing.T) {
	body := []byte("")
	len := messageLengthPrefixFormatter(len(body))
	payload := fmt.Sprintf("%s%s", len, body)
	r := strings.NewReader(payload)
	_, err := read(r, 100)

	if err == nil {
		t.Fatal("read must return error as the payload was empty")
	}
}

func TestSendPayload(t *testing.T) {

	b := bytes.NewBuffer(make([]byte, 0))
	msg := "THIS IS A TEST"

	err := write(b, []byte(msg))

	if err != nil {
		t.Fatalf("could not write to io.Writer, err:%s", err)
	}

	result := b.Bytes()

	if len(result) < 4+len(msg) {
		t.Fatalf("invalid size of buffer")
	}

	payload := string(result[4:])

	if payload != msg {
		t.Fatalf("the input and output don't match output: %s, input: %s", payload, msg)
	}

}

func createRandomString(n int) string {
	letterRunes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
