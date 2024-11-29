package main

import (
	"log"

	"github.com/Supriyo-455/CloudFileExplorer/p2p"
)

type FileServerOpts struct {
	ListenAddr  string
	StorageRoot string
	PathTransformFunc
	Transport p2p.Transport
}

type FileServer struct {
	FileServerOpts

	store  *Store
	quitch chan struct{}
}

func NewFileServer(opts FileServerOpts) *FileServer {
	storeOpts := StoreOpts{
		Root:              opts.StorageRoot,
		PathTransformFunc: opts.PathTransformFunc,
	}
	return &FileServer{
		FileServerOpts: opts,
		store:          NewStore(storeOpts),
		quitch:         make(chan struct{}),
	}
}

func (fs *FileServer) Start() error {
	if err := fs.Transport.ListenAndAccept(); err != nil {
		return err
	}
	return nil
}

func (fs *FileServer) loop() {
	for {
		select {
		case msg := <-fs.Transport.Consume():
			log.Println(msg)
		case <-fs.quitch:
			return
		}
	}
}

func (fs *FileServer) Stop() {
	close(fs.quitch)
}
