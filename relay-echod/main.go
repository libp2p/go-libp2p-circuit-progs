package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	libp2p "github.com/libp2p/go-libp2p"
	inet "github.com/libp2p/go-libp2p-net"
	pstore "github.com/libp2p/go-libp2p-peerstore"
	ma "github.com/multiformats/go-multiaddr"
)

const Proto = "/relay/test/echo"

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <relay-address>\n", os.Args[0])
		os.Exit(1)
	}

	raddr, err := ma.NewMultiaddr(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	rinfo, err := pstore.InfoFromP2pAddr(raddr)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	host, err := libp2p.New(
		ctx,
		libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/0/ws"),
		libp2p.EnableRelay(),
	)
	if err != nil {
		log.Fatal(err)
	}

	host.SetStreamHandler(Proto, handleStream)

	rctx, cancel := context.WithTimeout(ctx, time.Second)
	err = host.Connect(rctx, *rinfo)
	if err != nil {
		log.Fatal(err)
	}
	cancel()

	fmt.Printf("Listening at /p2p-circuit/ipfs/%s\n", host.ID().Pretty())
	select {}
}

func handleStream(s inet.Stream) {
	log.Printf("New echo stream from %s", s.Conn().RemoteMultiaddr().String())
	defer s.Close()
	count, err := io.Copy(s, s)
	if err != nil {
		log.Printf("Error echoing: %s", err.Error())
	}

	log.Printf("echoed %d bytes", count)
}
