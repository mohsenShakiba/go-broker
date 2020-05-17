package util

import (
	"bytes"
	"fmt"
	"testing"
)

func formatter(s string) []byte {
	prefix := []byte(fmt.Sprintf("%04d", len(s)+4))
	return append(prefix, []byte(s)...)
}

func TestReadSmall(t *testing.T) {
	msg := "THIS IS A TEST"
	payload := formatter(msg)
	r := bytes.NewReader(payload)
	result, err := Read(r, 100)

	if err != nil {
		t.Fatalf("error in reading, input: %s, error: %s", msg, err)
	}

	if string(result) != string(msg) {
		t.Fatalf("input and output mismatch while reading, input %s, output: %s", msg, result)
	}
}

func TestReadLarge(t *testing.T) {
	msg := "THIS IS A TEST"
	payload := formatter(msg)
	r := bytes.NewReader(payload)
	result, err := Read(r, 5)

	if err != nil {
		t.Fatalf("error in reading, input: %s, error: %s", msg, err)
	}

	if string(result) != msg {
		t.Fatalf("input and output mismatch while reading, input %s, output: %s", msg, result)
	}
}

func TestReadInvalid(t *testing.T) {
	msg := ""
	payload := formatter(msg)
	r := bytes.NewReader(payload)
	_, err := Read(r, 100)

	if err == nil {
		t.Fatal("Read must return error as the payload was empty")
	}
}
