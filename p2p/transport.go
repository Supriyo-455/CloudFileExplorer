package p2p

// NOTE: Peer represents the remote node
type Peer interface {
	Close() error
}

// NOTE: Transport handles communication between the nodes
// in the network. This can be of the form (TCP, UDP, Websockets, ...)
type Transport interface {
	ListenAndAccept() error

	// Consume implements the transport interface, which will return read-only channel
	// for reading the incoming messages received from another peer in the network.
	Consume() <-chan RPC
}
