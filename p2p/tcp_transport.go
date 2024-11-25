package p2p

import (
	"fmt"
	"net"
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
	OnPeer        func(Peer) error
}

type TCPTransport struct {
	TCPTransportOpts TCPTransportOps
	listener         net.Listener
	rpcch            chan RPC
}

func NewTCPTransport(opts TCPTransportOps) *TCPTransport {
	return &TCPTransport{
		TCPTransportOpts: opts,
		rpcch:            make(chan RPC),
	}
}

func (t *TCPTransport) Consume() <-chan RPC {
	return t.rpcch
}

func (p *TCPPeer) Close() error {
	return p.conn.Close()
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
	var err error

	defer func() {
		fmt.Printf("dropping peer connection: %s\n", err)
		conn.Close()
	}()

	peer := NewTCPPeer(conn, true)
	tcpOpts := t.TCPTransportOpts

	if err = tcpOpts.HandShakeFunc(peer); err != nil {
		return
	}

	if tcpOpts.OnPeer != nil {
		if err = tcpOpts.OnPeer(peer); err != nil {
			return
		}
	}

	// Read loop
	errCount := 0
	rpc := RPC{}
	for {
		if err = tcpOpts.Decoder.Decode(conn, &rpc); err != nil {
			errCount += 1
			if errCount == 5 {
				return
			}
			fmt.Printf("TCP Error: %s\n", err)
			continue
		}
		rpc.From = conn.RemoteAddr()
		t.rpcch <- rpc
	}
}
