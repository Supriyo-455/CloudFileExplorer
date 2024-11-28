package main

import (
	"bytes"
	"io"
	"testing"
)

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
	opts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	}
	s := NewStore(opts)
	key := "specialPicture"
	content := []byte("Hello Hunny Bunny!!")
	data := bytes.NewReader(content)

	if err := s.WriteStream(key, data); err != nil {
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
