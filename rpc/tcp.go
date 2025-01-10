package rpc

import (
	"fmt"
	"net"
)

type TCPPeer struct {
	conn     net.Conn
	outbound bool
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		conn:     conn,
		outbound: outbound,
	}
}

func (p *TCPPeer) Close() error {
	return p.conn.Close()
}

type TCPOpts struct {
	ListenAddress string
	Handshake     HandShake
	Decoder       Decoder
	PeerAttached  func(Peer) error
}

type TCPTransport struct {
	TCPOpts
	listener net.Listener
	rpcch    chan RPC
}

func NewTCPTransport(opts TCPOpts) *TCPTransport {
	return &TCPTransport{
		TCPOpts: opts,
		rpcch:   make(chan RPC),
	}
}

// impl Transport interface
// R.O. channel
func (t *TCPTransport) Consume() <-chan RPC {
	return t.rpcch
}

func (t *TCPTransport) ListenAndAccept() error {
	var err error

	t.listener, err = net.Listen("tcp", t.ListenAddress)
	if err != nil {
		return err
	}
	fmt.Println("Server started!")
	go t.acceptLoop()
	return nil
}

func (t *TCPTransport) acceptLoop() {
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			fmt.Printf("Failed to accept connection: %s\n ", err)
		}
		fmt.Printf("new incoming connection: %+v\n", conn)
		go t.handleConn(conn)
	}
}

func (t *TCPTransport) handleConn(conn net.Conn) {
	var err error

	defer func() {
		fmt.Printf("Dropping node conn: %s\n", err)
		conn.Close()
	}()

	p := NewTCPPeer(conn, true)

	//init handshake first
	if err = t.Handshake(p); err != nil {
		return
	}
	//then check that a peer is attached
	if t.PeerAttached != nil {
		if err = t.PeerAttached(p); err != nil {
			return
		}
	}

	rpc := RPC{}
	for {
		err := t.Decoder.Decode(conn, &rpc)
		if err != nil {
			return
		}
		rpc.From = conn.RemoteAddr()
		t.rpcch <- rpc
	}
}
