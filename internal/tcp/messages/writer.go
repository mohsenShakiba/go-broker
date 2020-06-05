package messages

import (
	"bufio"
	"fmt"
	log "github.com/sirupsen/logrus"
)

func WriteToIO(msg *Message, w *bufio.Writer) bool {

	// create header
	header := fmt.Sprintf("%s %d\n", msg.Type, len(msg.Fields))

	// write header
	_, err := w.Write([]byte(header))

	// if err
	if err != nil {
		log.Errorf("could not write messages, err: %s", err)
		return false
	}

	// for each field
	for k, v := range msg.Fields {

		// write key
		_, err := w.Write([]byte(k))

		// if err
		if err != nil {
			log.Errorf("could not write messages, err: %s", err)
			return false
		}

		_, err = w.Write(spaceSeparator)

		// if err
		if err != nil {
			log.Errorf("could not write messages, err: %s", err)
			return false
		}

		// write length
		// +1 for \n
		length := len(v) + 1
		lengthStr := fmt.Sprintf("%d", length)

		_, err = w.Write([]byte(lengthStr))

		// if err
		if err != nil {
			log.Errorf("could not write messages, err: %s", err)
			return false
		}

		_, err = w.Write(spaceSeparator)

		// if err
		if err != nil {
			log.Errorf("could not write messages, err: %s", err)
			return false
		}

		_, err = w.Write(v)

		// if err
		if err != nil {
			log.Errorf("could not write messages, err: %s", err)
			return false
		}

		_, err = w.Write(newLineSeparator)

	}

	return true
}
