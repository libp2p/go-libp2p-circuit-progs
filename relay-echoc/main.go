package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	//	"io/ioutil"
	"log"
	"os"
	"time"

	circuit "github.com/libp2p/go-libp2p-circuit"
	crypto "github.com/libp2p/go-libp2p-crypto"
	metrics "github.com/libp2p/go-libp2p-metrics"
	peer "github.com/libp2p/go-libp2p-peer"
	pstore "github.com/libp2p/go-libp2p-peerstore"
	swarm "github.com/libp2p/go-libp2p-swarm"
	bhost "github.com/libp2p/go-libp2p/p2p/host/basic"
	ma "github.com/multiformats/go-multiaddr"
)

const Proto = "/relay/test/echo"

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s <echod-address> <msg>\n", os.Args[0])
		os.Exit(1)
	}

	paddr, err := ma.NewMultiaddr(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	pinfo, err := pstore.InfoFromP2pAddr(paddr)
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
		[]ma.Multiaddr{},
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

	rctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	err = host.Connect(rctx, pstore.PeerInfo{pinfo.ID, []ma.Multiaddr{paddr}})
	if err != nil {
		log.Fatal(err)
	}

	s, err := host.NewStream(rctx, pinfo.ID, Proto)
	if err != nil {
		log.Fatal(err)
	}

	msg := []byte(os.Args[2])

	s.Write(msg)

	data := make([]byte, len(msg))
	_, err = s.Read(data)
	if err != nil && err != io.EOF {
		log.Fatal(err)
	}

	s.Close()

	fmt.Printf("Peer says: %s\n", string(data))

	if bytes.Equal(data, msg) {
		fmt.Printf("OK\n")
	} else {
		fmt.Printf("ERROR\n")
		os.Exit(1)
	}
}
