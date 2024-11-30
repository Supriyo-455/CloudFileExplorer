package main

import (
	"bytes"
	"log"
	"time"

	"github.com/Supriyo-455/CloudFileExplorer/p2p"
)

func makeServer(listenAddr string, root string, nodes []string) *FileServer {
	tcpTransportOpts := p2p.TCPTransportOps{
		ListenAddr:    listenAddr,
		HandShakeFunc: p2p.NOPHandshakeFunc,
		Decoder:       p2p.DefaultDecoder{},
	}

	tcpTransport := p2p.NewTCPTransport(tcpTransportOpts)

	fsOpts := FileServerOpts{
		ListenAddr:        tcpTransportOpts.ListenAddr,
		StorageRoot:       root,
		PathTransformFunc: CASPathTransformFunc,
		Transport:         tcpTransport,
		BootstrapNodes:    nodes,
	}

	s := NewFileServer(fsOpts)
	tcpTransport.OnPeer = s.OnPeer

	return s
}

func main() {
	s1 := makeServer(":3000", "HmzNetwork", []string{})
	s2 := makeServer(":4000", "HmzNetwork", []string{":3000"})

	go func() {
		log.Fatal(s1.Start())
	}()

	go func() {
		log.Fatal(s2.Start())
	}()

	time.Sleep(time.Second * 2)

	data := bytes.NewReader([]byte("my big data file here!"))
	if err := s2.StoreData("myPrivateData", data); err != nil {
		log.Fatal(err)
	}

	select {}
}
