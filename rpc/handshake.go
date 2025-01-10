package rpc

type HandShake func(Peer) error

func NoOpHandshake(Peer) error { return nil }
