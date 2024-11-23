package p2p

import "net"

// NOTE: Holds any data that is sent over
// each transport between two nodes in the network
type Message struct {
	From    net.Addr
	Payload []byte
}
