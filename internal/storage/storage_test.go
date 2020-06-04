package storage

import (
	"os"
	"path"
	"testing"
	"time"
)

func TestStorage(t *testing.T) {

	basePath := path.Join("C:\\Users\\m.shakiba.PSZ021-PC\\Downloads\\_tmp")

	err := os.Mkdir(basePath, os.ModeDir)

	if err != nil && !os.IsExist(err) {
		t.Fatalf("failed to create temp folder")
		return
	}

	defer func() {

		time.Sleep(time.Second)

		err := os.RemoveAll(basePath)

		if err != nil {
			t.Fatalf("failed to delete path")
		}
	}()

	conf := StorageConfig{
		Path:           basePath,
		FileMaxSize:    100,
		FileNamePrefix: "test_file_",
	}

	s := New(conf)
	//defer s.Dispose()

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

	time.Sleep(time.Second)

	// recreating the storage
	s2 := New(conf)

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

	// make sure the entry is actually deleted

	time.Sleep(time.Second)

	s3 := New(conf)

	err = s3.Init()

	if err != nil {
		t.Fatalf("initializing the storage failed with error: %s", err)
	}

	err = s3.Delete(1)

	if err != nil {
		t.Fatalf("failed to delete the entry, error: %s", err)
	}

	_, err = s3.Read(id)

	if err == nil {
		t.Fatalf("the entry should be deleted, error: %s", err)
	}

}
