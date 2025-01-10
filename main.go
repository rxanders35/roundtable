package main

// 2:15:03 in video
import (
	"fmt"
	"log"

	"github.com/rxanders35/roundtable/rpc"
)

func Attached(p rpc.Peer) error {
	p.Close()
	return nil
}

func main() {
	opts := rpc.TCPOpts{
		ListenAddress: ":3000",
		Handshake:     rpc.NoOpHandshake,
		Decoder:       rpc.DefaultDecoder{},
		PeerAttached:  Attached,
	}
	tr := rpc.NewTCPTransport(opts)

	go func() {
		for {
			msg := <-tr.Consume()
			fmt.Printf("Message %+v\n", msg)
		}
	}()
	if err := tr.ListenAndAccept(); err != nil {
		log.Fatal(err)
	}
	select {}
}
