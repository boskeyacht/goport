package main

import (
	"context"
	"database/sql"
	"fmt"
	"goport/config"
	"goport/db"
	"goport/listener"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/libp2p/go-libp2p"
	kad "github.com/libp2p/go-libp2p-kad-dht"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	cfg "github.com/libp2p/go-libp2p/config"
	multi "github.com/multiformats/go-multiaddr"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"
)

type Node struct {
	Bun      string // prisma client
	GraphQl  string // graphql server
	Ingestor string // ingestor
	Gossip   string // gossip
}

type NodeConfig struct{}

func (n *Node) New() {}

func (n *Node) Start() {

	// Open MySQL connection
	s, err := sql.Open("mysql", fmt.Sprintf("root:%s/%s", config.DB_NAME, config.DB_PASSWORD))
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err.Error())
		return
	}

	db := &db.SQLWrapper{
		DB: bun.NewDB(s, mysqldialect.New()),
	}

	// Create a new seaport contract listiner
	sl, err := listener.New()
	if err != nil {
		log.Fatalf("Failed to create new SeaportListener: %v", err.Error())
		return
	}

	// Create a new libp2p host
	lp, err := libp2p.New(initP2P())
	if err != nil {
		log.Fatalf("Failed to create new libp2p host: %v", err.Error())
		return
	}

	// Start the seaport listener
	sl.Start(db, &lp)

	// Create a new DHT
	dht := kad.NewDHT(context.Background(), lp, nil)

	// Create a new pubsub
	_, err = pubsub.NewGossipSub(context.Background(), lp)
	if err != nil {
		log.Fatalf("Failed to create gossipsub: %v", err.Error())
		return
	}

	err = dht.Bootstrap(context.Background())
	if err != nil {
		log.Fatalf("Failed to bootstrap DHT: %v", err.Error())
		return
	}

}

// Libp2p config function
func initP2P() cfg.Option {
	return func(c *cfg.Config) error {
		m, err := multi.NewMultiaddr(fmt.Sprintf("ip4/%s/tcp/%s", config.HOST_NAME, config.HOST_PORT))
		if err != nil {
			log.Printf("Invalid multiaddr: %v", err.Error())
			return err
		}

		c.DialTimeout = time.Duration(1000000000)
		c.ListenAddrs = []multi.Multiaddr{m} // c.apply

		return nil
	}
}
