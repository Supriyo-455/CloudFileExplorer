package main

import (
	"bytes"
	"encoding/gob"
	"io"
	"log"
	"sync"

	"github.com/Supriyo-455/CloudFileExplorer/p2p"
)

type FileServerOpts struct {
	ListenAddr     string
	StorageRoot    string
	Transport      p2p.Transport
	BootstrapNodes []string

	PathTransformFunc
}

type FileServer struct {
	peerLock sync.Mutex
	peers    map[string]p2p.Peer
	store    *Store
	quitch   chan struct{}

	FileServerOpts
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
		peers:          make(map[string]p2p.Peer),
	}
}

type RPC struct {
	Payload any
}

func (fs *FileServer) StoreData(key string, r io.Reader) error {
	// 1. Store this file to the disk
	// 2. Broadcast this file to all known peers in the network

	buf := new(bytes.Buffer)
	rpc := RPC{
		Payload: []byte("storageKey"),
	}
	if err := gob.NewEncoder(buf).Encode(rpc); err != nil {
		return err
	}

	for _, peer := range fs.peers {
		if err := peer.Send(buf.Bytes()); err != nil {
			return err
		}
	}

	return nil
}

func (fs *FileServer) OnPeer(peer p2p.Peer) error {
	fs.peerLock.Lock()
	defer fs.peerLock.Unlock()

	fs.peers[peer.RemoteAddr().String()] = peer

	log.Printf("connected with remote peer: %s\n", peer.RemoteAddr().String())

	return nil
}

func (fs *FileServer) bootstrapNetwork() {
	for _, addr := range fs.BootstrapNodes {
		go func(addr string) {
			log.Println("attempting to connect with remote: ", addr)
			if err := fs.Transport.Dial(addr); err != nil {
				log.Println("dial error: ", err)
			}
		}(addr)
	}
}

func (fs *FileServer) Start() error {
	if err := fs.Transport.ListenAndAccept(); err != nil {
		return err
	}

	fs.bootstrapNetwork()

	fs.loop()

	return nil
}

func (fs *FileServer) loop() {

	defer func() {
		log.Println("stopping the file server..")
		fs.Transport.Close()
	}()

	for {
		select {
		case rpcTransport := <-fs.Transport.Consume():
			var rpc RPC
			if err := gob.NewDecoder(bytes.NewReader(rpcTransport.Payload)).Decode(&rpc); err != nil {
				log.Fatal(err)
			}
			log.Println("received : ", string(rpc.Payload.([]byte)))

		case <-fs.quitch:
			return
		}
	}
}

func (fs *FileServer) Stop() {
	close(fs.quitch)
}
