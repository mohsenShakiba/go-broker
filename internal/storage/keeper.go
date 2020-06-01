package storage

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type keeper struct {
	basePath       string
	fileNamePrefix string
	maxFileSize    int64
	pages          []*page
	currentPage    *page
}

func newKeeper(cnf StorageConfig) *keeper {
	return &keeper{
		basePath:       cnf.Path,
		fileNamePrefix: cnf.FIleNamePrefix,
		maxFileSize:    cnf.FileMaxSize,
		pages:          make([]*page, 0),
		currentPage:    nil,
	}
}

func (k *keeper) init() error {
	files, err := ioutil.ReadDir(k.basePath)

	if err != nil {
		return err
	}

	for _, f := range files {
		if strings.Contains(f.Name(), k.fileNamePrefix) {

		}
	}
}

func (k *keeper) getHandler() (*os.File, error) {
	if

}
