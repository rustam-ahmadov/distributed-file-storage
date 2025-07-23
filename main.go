package main

import (
	"github.com/rustam-ahmadov/distributed-file-storage/p2p"
	"log/slog"
)

func main() {

	t := p2p.NewTCPTransport(":8080",
		p2p.WithDecoder(p2p.DefaultDecoder{}),
		p2p.WithHandShaker(p2p.NOPHandshakeFunc),
	)

	if err := t.ListenAndAccept(); err != nil {
		slog.Error("error from listen: ", err)
	}

	select {}
}
