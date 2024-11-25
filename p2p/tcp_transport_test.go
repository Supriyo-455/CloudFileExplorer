package p2p

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTCPTransport(t *testing.T) {
	opts := TCPTransportOps{
		ListenAddr:    ":3000",
		HandShakeFunc: NOPHandshakeFunc,
		Decoder:       DefaultDecoder{},
	}
	tr := NewTCPTransport(opts)

	assert.Equal(t, tr.TCPTransportOpts.ListenAddr, ":3000")
	assert.Nil(t, tr.ListenAndAccept())
}
