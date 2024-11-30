package p2p

import "net"

// NOTE: Peer represents the remote node
type Peer interface {
	// NOTE: conn is the underlying conn of the peer
	net.Conn
	Send([]byte) error
}

// NOTE: Transport handles communication between the nodes
// in the network. This can be of the form (TCP, UDP, Websockets, ...)
type Transport interface {
	Dial(string) error
	ListenAndAccept() error
	Consume() <-chan RPC
	Close() error
}
