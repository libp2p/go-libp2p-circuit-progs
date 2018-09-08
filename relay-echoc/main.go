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

	libp2p "github.com/libp2p/go-libp2p"
	pstore "github.com/libp2p/go-libp2p-peerstore"
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
	ctx := context.Background()
	host, err := libp2p.New(
		ctx,
		libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/0/ws"),
		libp2p.EnableRelay(),
	)
	if err != nil {
		log.Fatal(err)
	}

	rctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	err = host.Connect(rctx, pstore.PeerInfo{ID: pinfo.ID, Addrs: []ma.Multiaddr{paddr}})
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
