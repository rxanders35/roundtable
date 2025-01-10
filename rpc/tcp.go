package rpc

import (
	"fmt"
	"net"
	"sync"
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

type TCPOpts struct {
	ListenAddress string
	Handshake     HandShake
	Decoder       Decoder
}

type TCPTransport struct {
	TCPOpts
	listener net.Listener
	rpcch    chan RPC
	mu       sync.RWMutex
	peers    map[net.Addr]Peer
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
	p := NewTCPPeer(conn, true)
	if err := t.Handshake(p); err != nil {
		fmt.Printf("Error during handshake: %s\n", err)
		conn.Close()
		return
	}

	rpc := &RPC{}
	for {
		if err := t.Decoder.Decode(conn, rpc); err != nil {
			fmt.Printf("TCP error: %s\n", err)
			continue
		}
		rpc.From = conn.RemoteAddr()
		fmt.Println("from: ", rpc)
	}
}
