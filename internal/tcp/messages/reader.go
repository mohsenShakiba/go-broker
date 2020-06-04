package messages

import (
	"bufio"
	"bytes"
	log "github.com/sirupsen/logrus"
	"strconv"
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
		log.Errorf("messages part count is invalid, msg: %s", string(header))
		return nil, false
	}

	// messages parts
	msgType := string(headerParts[0])

	msgPartCountStr := string(headerParts[1])
	msgPartCount, err := strconv.Atoi(msgPartCountStr)

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
		key = key[:len(key)-1]

		// if error
		if err != nil {
			log.Errorf("could not read from socket, err: %s", err)
			return nil, false
		}

		// read space 2
		length, err := r.ReadSlice(' ')
		length = length[:len(length)-1]

		// if error
		if err != nil {
			log.Errorf("could not read from socket, err: %s", err)
			return nil, false
		}

		// read length
		msgPartCountStr := string(length)
		payloadLength, err := strconv.Atoi(msgPartCountStr)

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

	// get msg id
	msgId, ok := msg.Fields["msgId"]

	if !ok {
		log.Errorf("message doesn't contain a messages id, discarding")
		return nil, false
	}

	msg.MsgId = string(msgId)

	return msg, true
}
