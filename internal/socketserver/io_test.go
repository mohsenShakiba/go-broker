package socketserver

import (
	"bytes"
	"testing"
)

func TestReadSmall(t *testing.T) {
	msg := "THIS IS A TEST"
	payload := formatStr(msg)
	r := bytes.NewReader(payload)
	result, err := read(r, 100)

	if err != nil {
		t.Fatalf("error in reading, input: %s, error: %s", msg, err)
	}

	if string(result) != string(msg) {
		t.Fatalf("input and output mismatch while reading, input %s, output: %s", msg, result)
	}
}

func TestReadLarge(t *testing.T) {
	msg := "THIS IS A TEST"
	payload := formatStr(msg)
	r := bytes.NewReader(payload)
	result, err := read(r, 5)

	if err != nil {
		t.Fatalf("error in reading, input: %s, error: %s", msg, err)
	}

	if string(result) != msg {
		t.Fatalf("input and output mismatch while reading, input %s, output: %s", msg, result)
	}
}

func TestReadInvalid(t *testing.T) {
	msg := ""
	payload := formatStr(msg)
	r := bytes.NewReader(payload)
	_, err := read(r, 100)

	if err == nil {
		t.Fatal("read must return error as the payload was empty")
	}
}
