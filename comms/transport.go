package comms

// remote nodes
type Peer interface {
}

// node comms
type Transport interface {
	ListenAndAccept() error
}
