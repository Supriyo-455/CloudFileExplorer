package p2p

import (
	"errors"
	"log"
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

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		conn:     conn,
		outbound: outbound,
	}
}

// Close() implements peer interface
func (p *TCPPeer) Close() error {
	return p.conn.Close()
}

// RemoteAddr() implements peer interface and will return
// the remote address of its the underlying connection.
func (p *TCPPeer) RemoteAddr() net.Addr {
	return p.conn.RemoteAddr()
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

// Close() implements transport interface
func (t *TCPTransport) Close() error {
	return t.listener.Close()
}

// ListenAndAccept() implements transport interface
func (t *TCPTransport) ListenAndAccept() error {
	var err error

	t.listener, err = net.Listen("tcp", t.TCPTransportOpts.ListenAddr)
	if err != nil {
		return err
	}

	go t.startAcceptLoop()

	log.Printf("tcp transport listening on port: %s\n", t.TCPTransportOpts.ListenAddr)

	return nil
}

func (t *TCPTransport) startAcceptLoop() {
	for {
		conn, err := t.listener.Accept()

		if errors.Is(err, net.ErrClosed) {
			return
		}

		if err != nil {
			log.Printf("TCP accept error: %s\n", err)
		}

		log.Printf("new incoming connection %+v\n", conn)
		go t.handleConn(conn, false)
	}
}

func (t *TCPTransport) handleConn(conn net.Conn, outbound bool) {
	var err error

	defer func() {
		log.Printf("dropping peer connection: %s\n", err)
		conn.Close()
	}()

	peer := NewTCPPeer(conn, outbound)
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
			log.Printf("TCP Error: %s\n", err)
			continue
		}
		rpc.From = conn.RemoteAddr()
		t.rpcch <- rpc
	}
}

// Dial() Implements the transport interface
func (t *TCPTransport) Dial(addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}

	go t.handleConn(conn, true)

	return nil
}
