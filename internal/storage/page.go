package storage

import "os"

type page struct {
	path          string
	offset        int64
	activeEntries int
	fh            *os.File
}

// init will create the file if it doesn't exist,
func (p *page) createIfNeeded() error {
	if _, err := os.Stat(p.path); os.IsNotExist(err) {
		_, err = os.Create(p.path)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *page) fileHandler() (*os.File, error) {
	if p.fh != nil {
		return p.fh, nil
	}

	fh, err := os.Open(p.path)

	if err != nil {
		return nil, err
	}

	p.fh = fh

	return fh, nil
}

func (p *page) close() {
	if p.fh != nil {
		_ = p.fh.Close()
	}
}
