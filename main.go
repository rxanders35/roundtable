package main

// 1:00:33 in video
import (
	"log"

	"github.com/rxanders35/roundtable/rpc"
)

func main() {
	opts := rpc.TCPOpts{
		ListenAddress: ":3000",
		Handshake:     rpc.NoOpHandshake,
		Decoder:       rpc.DefaultDecoder{},
	}
	tr := rpc.NewTCPTransport(opts)
	if err := tr.ListenAndAccept(); err != nil {
		log.Fatal(err)
	}
	select {}
}
