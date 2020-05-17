package util

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"io"
	"strconv"
)

func Read(i io.Reader, size int) ([]byte, error) {

	msg := make([]byte, 0, size)
	buf := make([]byte, size)
	remainingBytesToRead := 0
	bytesToReadInit := false

	for {
		bytesRead, err := i.Read(buf)

		if !bytesToReadInit && bytesRead > 4 {
			sizeStr := string(buf[:4])
			remainingBytesToRead, err = strconv.Atoi(sizeStr)

			if err != nil {
				log.Errorf("the first 4 bytes of message isn't int, first four bytes %s", sizeStr)
				break
			}

			bytesToReadInit = true

		}

		msg = append(msg, buf[:bytesRead]...)

		remainingBytesToRead -= bytesRead

		if remainingBytesToRead < 0 {
			return nil, errors.New("invalid message was sent from client")
		}

		if remainingBytesToRead == 0 {
			break
		}
	}

	if !bytesToReadInit {
		return nil, errors.New("couldn't Read due to invalid message format")
	}

	log.Infof("received message from socket, msg: %s", msg)

	lStr := string(msg[:4])
	l, err := strconv.Atoi(lStr)

	if err != nil {
		log.Fatalf("message length is invalid %s", lStr)
		return nil, err
	}

	return msg[4:l], nil
}
