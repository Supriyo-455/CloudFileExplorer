package p2p

// NOTE: Peer represents the remote node
type Peer interface {
}

// NOTE: Transport handles communication between the nodes
// in the network. This can be of the form (TCP, UDP, Websockets, ...)
type Transport interface {
	ListenAndAccept() error
}
