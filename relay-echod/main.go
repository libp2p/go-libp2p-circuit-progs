package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	circuit "github.com/libp2p/go-libp2p-circuit"
	crypto "github.com/libp2p/go-libp2p-crypto"
	metrics "github.com/libp2p/go-libp2p-metrics"
	inet "github.com/libp2p/go-libp2p-net"
	peer "github.com/libp2p/go-libp2p-peer"
	pstore "github.com/libp2p/go-libp2p-peerstore"
	swarm "github.com/libp2p/go-libp2p-swarm"
	bhost "github.com/libp2p/go-libp2p/p2p/host/basic"
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

	// need a binding to be able to dial ws addresses
	wsaddr, err := ma.NewMultiaddr("/ip4/0.0.0.0/tcp/0/ws")
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
		[]ma.Multiaddr{wsaddr},
		id,
		ps,
		metrics.NewBandwidthCounter(),
	)
	if err != nil {
		log.Fatal(err)
	}

	host := bhost.New(netw)
	err = circuit.AddRelayTransport(ctx, host)
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

	fmt.Printf("Listening at /p2p-circuit/ipfs/%s\n", id.Pretty())
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
