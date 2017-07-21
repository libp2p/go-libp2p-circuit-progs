go-libp2p-circuit-progs
======================

[![](https://img.shields.io/badge/made%20by-Protocol%20Labs-blue.svg?style=flat-square)](http://ipn.io)
[![](https://img.shields.io/badge/project-IPFS-blue.svg?style=flat-square)](http://libp2p.io/)
[![](https://img.shields.io/badge/freenode-%23ipfs-blue.svg?style=flat-square)](http://webchat.freenode.net/?channels=%23ipfs)

Simple (interop) testing programs for go-libp2p-circuit:
- relay, a Hop relay.
- relay-echod: an echoing daemon reachable through relay.
- relay-echoc: a client for an echo daemin reachable through relay.

## Install

```sh
go get github.com/libp2p/go-libp2p-circuit-progs/...
```

## Usage

Start a relay:
```sh
$ relay
Relay addresses:
/ip4/127.0.0.1/tcp/9001/ipfs/QmRhDVAfdaBw7zpk5jn4L4BcBJDH1DZrKBHENqUwAvjkRe
/ip4/10.0.1.30/tcp/9001/ipfs/QmRhDVAfdaBw7zpk5jn4L4BcBJDH1DZrKBHENqUwAvjkRe
/ip6/::1/tcp/9001/ipfs/QmRhDVAfdaBw7zpk5jn4L4BcBJDH1DZrKBHENqUwAvjkRe

```

Start an echo server reachable through relay:
```sh
$ relay-echod /ip4/127.0.0.1/tcp/9001/ipfs/QmRhDVAfdaBw7zpk5jn4L4BcBJDH1DZrKBHENqUwAvjkRe
Listening at /p2p-circuit/ipfs/QmUB1XXEDrXZqXGTDBNtcniUNv9ZAE8bEZnqCoYVXyV6Bc
```

Use the echo client to talk to an echo server:
```sh
$ relay-echoc /ip4/127.0.0.1/tcp/9001/ipfs/QmRhDVAfdaBw7zpk5jn4L4BcBJDH1DZrKBHENqUwAvjkRe/p2p-circuit/ipfs/QmUB1XXEDrXZqXGTDBNtcniUNv9ZAE8bEZnqCoYVXyV6Bc "hello world"
Peer says: hello world
OK

```

## Contribute

PRs are welcome!

## License

MIT (c) vyzo
