package p2p

import (
	"fmt"
	"net"
	"sync"
)

// NOTE: Represents the remote node over
// a tcp-established connection
type TCPPeer struct {
	// NOTE: conn is the underlying conn of the peer
	conn net.Conn

	// NOTE: if we dial and retrieve a conn => outbound == true
	// if we accept and retrieve a conn => outbound == false
	outbound bool
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		conn:     conn,
		outbound: outbound,
	}
}

type TCPTransportOps struct {
	ListenAddr    string
	HandShakeFunc HandShakeFunc
	Decoder       Decoder
}

type TCPTransport struct {
	TCPTransportOpts TCPTransportOps
	listener         net.Listener

	mu    sync.RWMutex
	peers map[net.Addr]Peer
}

func NewTCPTransport(opts TCPTransportOps) *TCPTransport {
	return &TCPTransport{
		TCPTransportOpts: opts,
	}
}

func (t *TCPTransport) ListenAndAccept() error {
	var err error

	t.listener, err = net.Listen("tcp", t.TCPTransportOpts.ListenAddr)
	if err != nil {
		return err
	}

	go t.startAcceptLoop()

	return nil
}

func (t *TCPTransport) startAcceptLoop() {
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			fmt.Printf("TCP accept error: %s\n", err)
		}
		fmt.Printf("new incoming connection %+v\n", conn)
		go t.handleConn(conn)
	}
}

func (t *TCPTransport) handleConn(conn net.Conn) {
	peer := NewTCPPeer(conn, true)
	tcpOpts := t.TCPTransportOpts

	if err := tcpOpts.HandShakeFunc(peer); err != nil {
		conn.Close()
		fmt.Printf("TCP handshake error: %s\n", err)
		return
	}

	// Read loop
	errCount := 0
	msg := &Message{}
	for {
		if err := tcpOpts.Decoder.Decode(conn, msg); err != nil {
			errCount += 1
			if errCount == 5 {
				break
			}
			fmt.Printf("TCP Error: %s\n", err)
			continue
		}
		fmt.Printf("message: %+v\n", msg)
	}
}
