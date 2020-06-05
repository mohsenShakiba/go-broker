package storage

import (
	"os"
	"path"
	"testing"
	"time"
)

func TestStorageFile(t *testing.T) {
	basePath := path.Join("C:\\Users\\user\\Desktop")

	//fh1, err := os.OpenFile(basePath, os.O_RDWR, 0777)
	//
	//b := make([]byte, 20)
	//fh1.Read(b)
	//
	//e := fromBinary(b)
	//
	//fh, err := os.OpenFile(basePath, os.O_RDWR, 0777)
	//
	//if err != nil {
	//	t.Fatalf("failed to open the file, err: %s", err)
	//}
	//
	//if e == nil {
	//	t.Fatalf("failed to read binary")
	//}
	//
	//_, err = fh.WriteAt([]byte("test"),20)
	//
	//if err != nil {
	//	t.Fatalf("failed to write to file, err: %s", err)
	//}

	conf := StorageConfig{
		Path:           basePath,
		FileMaxSize:    100,
		FileNamePrefix: "tf_",
	}

	s := New(conf)
	//defer s.Dispose()

	err := s.Init()
	if err != nil {
		t.Fatal(err)
	}

	s.Write(1, []byte("test"))
	s.Write(2, []byte("test2"))
	s.Write(3, []byte("test3"))

	e, err := s.Read(1)

	if err != nil {
		t.Fatalf("the entry wasn't found")
	}

	t.Logf("e has id %s", e)

	s.Delete(1)

}

func TestWrite(t *testing.T) {

}

func TestRead(t *testing.T) {

}

func TestDelete(t *testing.T) {

}

func TestStorage(t *testing.T) {

	basePath := path.Join("C:\\Users\\user\\Desktop\\tmp")

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
