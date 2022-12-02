package main

import (
	"fmt"
	"goport/config"
	"goport/node"
	"log"
	"sync"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/peerstore"
	multi "github.com/multiformats/go-multiaddr"
)

func main() {
	wg := &sync.WaitGroup{}

	log.Println("Starting node...")

	m, err := multi.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%s", config.HOST_NAME, config.HOST_PORT))
	if err != nil {
		log.Printf("Invalid multiaddr: %v", err.Error())
		return
	}

	// If you would like to connect to a specific peer, you can do so here by adding it to the peerstore.
	// To add it, first `Add()` the peer to a newly created peerstore variable of type Peerstore.
	// Then pass `libp2p.Peerstore(peerstore_var)` as an additional option to the libp2p constructor.
	//
	// Types for constructing a peer for a peerstore can be found below:
	// https://pkg.go.dev/time#Duration
	// https://pkg.go.dev/github.com/multiformats/go-multiaddr#Multiaddr
	// https://pkg.go.dev/github.com/libp2p/go-libp2p@v0.23.4/core/peer#ID
	var p peerstore.Peerstore
	n, err := node.New(libp2p.ListenAddrs(m), libp2p.Peerstore(p))
	if err != nil {
		log.Printf("Failed to create new node: %v", err.Error())
		return
	}

	if err := n.Start(wg); err != nil {
		log.Printf("Failed to start node: %v", err.Error())
		return
	}

	wg.Wait()
}
