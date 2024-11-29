package main

import (
	"bytes"
	"fmt"
	"io"
	"testing"
)

func newStore() *Store {
	opts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	}
	return NewStore(opts)
}

func tearDown(t *testing.T, s *Store) {
	if err := s.Clear(); err != nil {
		t.Error(err)
	}
}

func TestPathTransformFunc(t *testing.T) {
	key := "momsbestdish"
	pathKey := CASPathTransformFunc(key)
	pathName := pathKey.PathName
	expectedPathName := "89e4c/9a0d6/7ba03/85064/c3ff7/5ec0b/b3f27/af2b5"
	FileNameKey := "89e4c9a0d67ba0385064c3ff75ec0bb3f27af2b5"

	if pathName != expectedPathName {
		t.Errorf("Expected path name: %s\n Actual path name:%s\n", expectedPathName, pathName)
	}

	if pathKey.FileName != FileNameKey {
		t.Errorf("Expected FileName Key: %s\n Actual path key:%s\n", FileNameKey, pathKey.FileName)
	}
}

func TestStore(t *testing.T) {
	s := newStore()
	defer tearDown(t, s)

	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("foo_%d", i)
		content := []byte("Hello Hunny Bunny!!")
		data := bytes.NewReader(content)

		if err := s.Write(key, data); err != nil {
			t.Error(err)
		}

		r, err := s.Read(key)
		if err != nil {
			t.Error(err)
		}

		b, err := io.ReadAll(r)
		if err != nil {
			t.Error(err)
		}

		if string(b) != string(content) {
			t.Errorf("Expected: %s, Got: %s\n", string(content), string(b))
		}

		if err := s.Delete(key); err != nil {
			t.Error(err)
		}

		if s.HasKey(key) {
			t.Errorf("key [%s] should not exist after delete operation!", key)
		}
	}
}
