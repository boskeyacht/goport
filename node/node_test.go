package node

import (
	"fmt"
	_ "goport/config"
	"sync"
	"testing"

	"github.com/libp2p/go-libp2p"
	multi "github.com/multiformats/go-multiaddr"
)

func TestNewNode(t *testing.T) {
	m, err := multi.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%s", "0.0.0.0", "8080"))
	if err != nil {
		t.Fatalf("Invalid multiaddr: %v", err.Error())
		return
	}

	n, err := New(libp2p.ListenAddrs(m))
	if err != nil || n == nil {
		t.Fatalf("Failed to create new node: %v", err.Error())
	}

	wg := &sync.WaitGroup{}
	err = n.Start(wg)

	if err != nil {
		t.Fatalf("Failed to start node: %v", err.Error())
	}

	err = n.Host.Close()
	if err != nil {
		t.Fatalf("Failed to stop node: %v", err.Error())
	}
}
