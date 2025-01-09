package main

import (
	"log"

	"github.com/rxanders35/roundtable/comms"
)

func main() {
	opts := comms.TCPOpts{
		ListenAddress: ":3000",
		Handshake:     comms.deniedHandshake,
	}
	tr := comms.NewTCPTransport(":3000")
	if err := tr.ListenAndAccept(); err != nil {
		log.Fatal(err)
	}
	select {}
}
