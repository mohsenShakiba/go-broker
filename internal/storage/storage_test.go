package storage

import (
	"os"
	"path"
	"testing"
)

func TestStorage(t *testing.T) {

	basePath := path.Join("C:\\Users\\m.shakiba.PSZ021-PC\\Desktop", "_tmp")

	err := os.Mkdir(basePath, os.ModeDir)

	if err != nil && !os.IsExist(err) {
		t.Fatalf("failed to create temp folder")
		return
	}

	defer func() {
		err := os.RemoveAll(basePath)

		if err != nil {
			//t.Fatalf("failed to delete path")
		}
	}()

	conf := StorageConfig{
		Path:           basePath,
		FileMaxSize:    100,
		FileNamePrefix: "test_file_",
	}

	s := New(conf)
	defer s.Dispose()

	err = s.Init()

	if err != nil {
		t.Fatalf("initializing the storage failed with error: %s", err)
	}

	val := "TEST"
	var id int64 = 1

	err = s.Write(1, []byte(val))

	if err != nil {
		t.Fatalf("writing to storage failed with error: %s", err)
	}

	// recreating the storage
	s2 := New(conf)

	defer s2.Dispose()

	err = s2.Init()

	if err != nil {
		t.Fatalf("initializing the storage failed with error: %s", err)
	}

	res, err := s2.Read(id)

	if err != nil {
		t.Fatalf("failed to read data from storage, error: %s", err)
	}

	if string(res) != val {
		t.Fatalf("failed to read data from storage, %s != %s", string(res), val)
	}
}
