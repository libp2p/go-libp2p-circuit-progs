package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	circuit "github.com/libp2p/go-libp2p-circuit"
	crypto "github.com/libp2p/go-libp2p-crypto"
	metrics "github.com/libp2p/go-libp2p-metrics"
	peer "github.com/libp2p/go-libp2p-peer"
	pstore "github.com/libp2p/go-libp2p-peerstore"
	swarm "github.com/libp2p/go-libp2p-swarm"
	bhost "github.com/libp2p/go-libp2p/p2p/host/basic"
	ma "github.com/multiformats/go-multiaddr"
)

func main() {
	port := flag.Int("l", 9001, "Relay listen port")
	flag.Parse()

	if len(flag.Args()) != 0 {
		fmt.Fprintf(os.Stderr, "Usage: %s [options ...]\nOptions:\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	ip4addr, err := ma.NewMultiaddr(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", *port))
	if err != nil {
		log.Fatal(err)
	}

	ip6addr, err := ma.NewMultiaddr(fmt.Sprintf("/ip6/::/tcp/%d", *port))
	if err != nil {
		log.Fatal(err)
	}

	privk, pubk, err := crypto.GenerateKeyPair(crypto.RSA, 2048)
	if err != nil {
		log.Fatal(err)
	}

	id, err := peer.IDFromPrivateKey(privk)
	if err != nil {
		log.Fatal(err)
	}

	ps := pstore.NewPeerstore()
	ps.AddPrivKey(id, privk)
	ps.AddPubKey(id, pubk)

	ctx := context.Background()
	netw, err := swarm.NewNetwork(
		ctx,
		[]ma.Multiaddr{ip4addr, ip6addr},
		id,
		ps,
		metrics.NewBandwidthCounter(),
	)
	if err != nil {
		log.Fatal(err)
	}

	host := bhost.New(netw)
	err = circuit.AddRelayTransport(ctx, host, circuit.OptHop)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Relay addresses:\n")
	for _, addr := range host.Addrs() {
		_, err := addr.ValueForProtocol(circuit.P_CIRCUIT)
		if err == nil {
			continue
		}
		fmt.Printf("%s/ipfs/%s\n", addr.String(), id.Pretty())
	}

	select {}
}
