package p2p

// NOTE: Holds any data that is sent over
// each transport between two nodes in the network
type Message struct {
	Payload []byte
}
