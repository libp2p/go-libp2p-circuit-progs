package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	libp2p "github.com/libp2p/go-libp2p"
	circuit "github.com/libp2p/go-libp2p-circuit"
)

func main() {
	port := flag.Int("l", 9001, "Relay TCP listen port")
	wsport := flag.Int("ws", 9002, "Relay WS listen port")
	flag.Parse()

	if len(flag.Args()) != 0 {
		fmt.Fprintf(os.Stderr, "Usage: %s [options ...]\nOptions:\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	ctx := context.Background()
	host, err := libp2p.New(
		ctx,
		libp2p.ListenAddrStrings(
			fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", *port),
			fmt.Sprintf("/ip6/::/tcp/%d", *port),
			fmt.Sprintf("/ip4/0.0.0.0/tcp/%d/ws", *wsport),
		),
		libp2p.EnableRelay(circuit.OptHop),
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Relay addresses:\n")
	for _, addr := range host.Addrs() {
		_, err := addr.ValueForProtocol(circuit.P_CIRCUIT)
		if err == nil {
			continue
		}
		fmt.Printf("%s/ipfs/%s\n", addr.String(), host.ID().Pretty())
	}

	select {}
}
