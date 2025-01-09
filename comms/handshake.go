package comms

type HandShake func(Peer) error

func deniedHandshake(Peer) error { return nil }
