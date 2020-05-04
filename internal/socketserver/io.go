package socketserver

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"strconv"
)

func toByteWithLengthPrefix(body string) []byte {
	msg := fmt.Sprintf("%04d%s", len(body) + 4, body)
	return []byte(msg)
}

func read(i io.Reader, size int) ([]byte, error) {
	log.Infof("starting to received from socket")

	msg := make([]byte, size)
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

		log.Infof("read %d bytes from socket with data %s, remaining %d", bytesRead, string(buf[:bytesRead]), remainingBytesToRead - bytesRead)

		if remainingBytesToRead == 0 {
			break
		}
	}

	if !bytesToReadInit {
		return nil, errors.New("couldn't read due to invalid message format")
	}

	log.Infof("received message from socket, msg: %s", msg)

	len, err := strconv.Atoi(string(msg[:4]))

	if err != nil {
		log.Fatalf("message length is invalid %d", len)
		return nil, err
	}

	return msg[4:len], nil
}
