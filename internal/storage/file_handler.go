package storage

import (
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
)

type fileHandler struct {
	conf             StorageConfig
	lock             sync.RWMutex
	pages            map[int]*page
	currentPage      *page
	currentPageIndex int
}

func newHandler(conf StorageConfig) *fileHandler {
	return &fileHandler{
		conf:             conf,
		lock:             sync.RWMutex{},
		pages:            make(map[int]*page),
		currentPage:      nil,
		currentPageIndex: 0,
	}
}

func (fh *fileHandler) readAllEntries() ([]*entry, error) {
	fh.lock.Lock()
	defer fh.lock.Unlock()

	files, err := ioutil.ReadDir(fh.conf.Path)

	if err != nil {
		return nil, err
	}

	entries := make([]*entry, 0, 1024)

	for _, f := range files {
		if strings.Contains(f.Name(), fh.conf.FIleNamePrefix) {

			pageIndexStr := strings.Replace(f.Name(), fh.conf.FIleNamePrefix, "", 1)

			pageIndex, err := strconv.Atoi(pageIndexStr)

			if err != nil {
				return nil, err
			}

			if pageIndex > fh.currentPageIndex {
				fh.currentPageIndex = pageIndex
			}

			filePath := path.Join(fh.conf.Path, f.Name())
			h, err := os.Open(filePath)

			if err != nil {
				return nil, err
			}

			var off int64 = 0
			b := make([]byte, 0, 18)
			blen := int64(18)
			pageEntries := make([]*entry, 0, 1024)

			for {
				_, err := h.ReadAt(b, off)

				if err == io.EOF {
					break
				}

				if err != nil {
					return nil, err
				}

				entry := fromBinary(b)

				entry.offset = off
				off += blen + int64(entry.length)

				pageEntries = append(pageEntries, entry)
			}

			p := &page{
				path:          filePath,
				offset:        off,
				activeEntries: len(pageEntries),
				fh:            h,
			}

			entries = append(entries, pageEntries...)
			fh.pages[pageIndex] = p
			fh.currentPage = p
		}
	}

	return entries, nil
}

func (fh *fileHandler) readPayload(e *entry) ([]byte, error) {
	fh.lock.RLock()
	defer fh.lock.RUnlock()

	// read the file content
	handler, err := fh.currentPage.fileHandler()

	if err != nil {
		return nil, err
	}

	b := make([]byte, e.length)
	_, err = handler.ReadAt(b, e.offset+18)

	if err != nil {
		return nil, err
	}

	return b, nil
}

func (fh *fileHandler) write(id int64, payload []byte) (*entry, error) {
	p := fh.currentPage

	if fh.currentPage.offset >= fh.conf.FileMaxSize {

		fh.currentPageIndex += 1
		newPName := fmt.Sprintf("%s%d", fh.conf.FIleNamePrefix, fh.currentPageIndex)
		newPPath := path.Join(fh.conf.Path, newPName)

		p := &page{
			path:          newPPath,
			offset:        0,
			activeEntries: 0,
			fh:            nil,
		}

		err := p.createIfNeeded()

		if err != nil {
			return nil, err
		}

		fh.pages[fh.currentPageIndex] = p
		fh.currentPage = p
	}

	entry := &entry{
		deleted: 0,
		id:      id,
		length:  int64(len(payload)),
		offset:  p.offset,
		pageNo:  int16(fh.currentPageIndex),
	}

	h, err := p.fileHandler()

	if err != nil {
		return nil, err
	}

	b := toBinary(entry)

	_, _ = h.WriteAt(b, p.offset)
	p.offset += int64(len(b))

	_, _ = h.WriteAt(payload, p.offset)
	p.offset += int64(len(payload))

	return entry, nil
}

func (fh *fileHandler) delete(entry *entry) error {
	p := fh.pages[int(entry.pageNo)]

	h, err := p.fileHandler()

	if err != nil {
		return err
	}

	b := make([]byte, 2)
	binary.LittleEndian.PutUint16(b, 0)
	_, err = h.WriteAt(b, entry.offset)

	if err != nil {
		return err
	}

	return nil
}
