package message

import (
	"bufio"
	"bytes"
	log "github.com/sirupsen/logrus"
	"go-broker/internal/tcp/util"
)

var spaceSeparator = []byte(" ")

func ReadFromIO(r *bufio.Reader) (*Message, bool) {

	// read header
	header, err := r.ReadSlice('\n')

	// if error
	if err != nil {
		log.Errorf("could not read from socket, err: %s", err)
		return nil, false
	}

	// trim the \n
	header = header[:len(header)-1]

	// parse the header
	headerParts := bytes.Split(header, spaceSeparator)

	// check if length of header is valid
	if len(headerParts) != 2 {
		log.Errorf("message part count is invalid, msg: %s", string(header))
		return nil, false
	}

	// message parts
	msgType := string(headerParts[0])
	msgPartCount, err := util.ReadLength(headerParts[1])

	// if can't parse length
	if err != nil {
		log.Errorf("failed to parse part count, err: %s", err)
		return nil, false
	}

	msg := &Message{
		Type:   msgType,
		Fields: make(map[string][]byte),
	}

	// for each part
	for i := 0; i < msgPartCount; i++ {

		// read key
		key, err := r.ReadSlice(' ')

		// if error
		if err != nil {
			log.Errorf("could not read from socket, err: %s", err)
			return nil, false
		}

		// read space 2
		length, err := r.ReadSlice(' ')

		// if error
		if err != nil {
			log.Errorf("could not read from socket, err: %s", err)
			return nil, false
		}

		// read length
		payloadLength, err := util.ReadLength(length)

		// if can't parse length
		if err != nil {
			log.Errorf("failed to parse part length, err: %s", err)
			return nil, false
		}

		// payload
		payload := make([]byte, payloadLength)

		// read
		_, err = r.Read(payload)

		// if error
		if err != nil {
			log.Errorf("failed to read payload, err: %s", err)
			return nil, false
		}

		msg.Fields[string(key)] = payload
	}

	return msg, true
}
