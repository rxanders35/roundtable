package rpc

// remote nodes
type Peer interface {
	Close() error
}

// node comms
type Transport interface {
	ListenAndAccept() error
	Consume() <-chan RPC
}
