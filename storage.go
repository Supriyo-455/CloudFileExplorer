package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const DEFAULT_ROOT = "hmzNetwork"

// NOTE: Lots of types of transform funcs can be added, like git paths etc. etc..
type PathTransformFunc func(string) PathKey

type PathKey struct {
	PathName string
	FileName string
}

type StoreOpts struct {
	PathTransformFunc PathTransformFunc
	Root              string
}

var DefaultPathTransformFunc = func(key string) PathKey {
	return PathKey{key, key}
}

func (p PathKey) FullPath() string {
	return fmt.Sprintf("%s/%s", p.PathName, p.FileName)
}

func CASPathTransformFunc(key string) PathKey {
	hash := sha1.Sum([]byte(key)) // [20]byte => []byte => [:]
	hashStr := hex.EncodeToString(hash[:])

	blockSize := 5
	sliceLen := len(hashStr) / blockSize
	paths := make([]string, sliceLen)

	for i := 0; i < sliceLen; i++ {
		from, to := i*blockSize, (i*blockSize)+blockSize
		paths[i] = hashStr[from:to]
	}

	return PathKey{PathName: strings.Join(paths, "/"), FileName: hashStr}
}

type Store struct {
	StoreOpts
}

func NewStore(opts StoreOpts) *Store {
	if opts.PathTransformFunc == nil {
		opts.PathTransformFunc = DefaultPathTransformFunc
	}
	if len(opts.Root) == 0 {
		opts.Root = DEFAULT_ROOT
	}
	return &Store{
		StoreOpts: opts,
	}
}

func (s *Store) HasKey(key string) bool {
	pathKey := s.PathTransformFunc(key)
	fullPath := pathKey.FullPath()
	fullPathWithRoot := s.Root + "/" + fullPath

	_, err := os.Stat(fullPathWithRoot)
	return !errors.Is(err, os.ErrNotExist)
}

func (s *Store) Clear() error {
	return os.RemoveAll(s.Root)
}

func (s *Store) Delete(key string) error {
	pathKey := s.PathTransformFunc(key)
	fullPath := pathKey.FullPath()
	fullPathWithRoot := s.Root + "/" + fullPath

	defer func() {
		log.Printf("deleted [%s] from disk", pathKey.FileName)
	}()

	return os.Remove(fullPathWithRoot)
}

func (s *Store) Read(key string) (io.Reader, error) {
	f, err := s.readStream(key)
	if err != nil {
		return nil, nil
	}

	defer f.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, f)

	return buf, err
}

func (s *Store) Write(key string, r io.Reader) error {
	return s.writeStream(key, r)
}

func (s *Store) readStream(key string) (io.ReadCloser, error) {
	pathKey := s.PathTransformFunc(key)
	fullPath := pathKey.FullPath()
	fullPathWithRoot := s.Root + "/" + fullPath
	return os.Open(fullPathWithRoot)
}

func (s *Store) writeStream(key string, r io.Reader) error {
	pathKey := s.PathTransformFunc(key)
	pathName := pathKey.PathName
	pathNameWithRoot := s.Root + "/" + pathName
	if err := os.MkdirAll(pathNameWithRoot, os.ModePerm); err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	io.Copy(buf, r)

	fullPath := pathKey.FullPath()
	fullPathWithRoot := s.Root + "/" + fullPath

	file, err := os.Create(fullPathWithRoot)
	if err != nil {
		return err
	}

	n, err := io.Copy(file, buf)
	if err != nil {
		return err
	}

	log.Printf("written (%d) bytes to disk: %s\n", n, fullPathWithRoot)

	return nil
}
