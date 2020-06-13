package storage

import (
	log "github.com/sirupsen/logrus"
	"os"
)

type page struct {
	path          string
	offset        int64
	activeEntries int
	fh            *os.File
}

func (p *page) fileHandler() (*os.File, error) {
	return os.OpenFile(p.path, os.O_RDWR|os.O_CREATE, 0777)
}

func (p *page) close() {
	if p.fh != nil {
		err := p.fh.Close()

		if err != nil {
			log.Errorf("failed to close file handler, error: %s", err)
		}
	}
}
