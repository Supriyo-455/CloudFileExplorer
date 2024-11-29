package main

import (
	"log"

	"github.com/Supriyo-455/CloudFileExplorer/p2p"
)

func main() {
	tcpTransportOpts := p2p.TCPTransportOps{
		ListenAddr:    ":3000",
		HandShakeFunc: p2p.NOPHandshakeFunc,
		Decoder:       p2p.DefaultDecoder{},
		// TODO: OnPeer function
	}

	tcpTransport := p2p.NewTCPTransport(tcpTransportOpts)

	fsOpts := FileServerOpts{
		ListenAddr:        tcpTransportOpts.ListenAddr,
		StorageRoot:       "3000_network",
		PathTransformFunc: CASPathTransformFunc,
		Transport:         tcpTransport,
	}

	fs := NewFileServer(fsOpts)

	if err := fs.Start(); err != nil {
		log.Fatal(err)
	}

	select {}
}
