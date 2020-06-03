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
	if p.fh != nil {
		return p.fh, nil
	}

	fh, err := os.OpenFile(p.path, os.O_CREATE|os.O_RDWR, 0644)

	if err != nil {
		return nil, err
	}

	p.fh = fh

	return fh, nil
}

func (p *page) close() {
	if p.fh != nil {
		err := p.fh.Close()

		if err != nil {
			log.Errorf("failed to close file handler, error: %s", err)
		}
	}
}
