package node

import (
	"context"
	"database/sql"
	"fmt"
	"goport/config"
	"goport/db"
	"goport/listener"
	"log"
	"sync"

	"github.com/libp2p/go-libp2p"
	kad "github.com/libp2p/go-libp2p-kad-dht"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	cfg "github.com/libp2p/go-libp2p/config"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
)

type Node struct {
	Host host.Host
	db.SQLWrapper
}

type NodeConfig struct{}

// Create a new libp2p host
func New(option ...cfg.Option) (*Node, error) {
	lp, err := libp2p.New(option...)
	if err != nil {
		log.Fatalf("Failed to create new libp2p host: %v", err.Error())
		return nil, err
	}

	return &Node{
		Host: lp,
	}, nil
}

// Start the node
func (n *Node) Start(wg *sync.WaitGroup) error {
	// Open MySQL connection
	s, err := sql.Open(sqliteshim.ShimName, fmt.Sprintf("file:%s:", config.DB_NAME))
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err.Error())
		return err
	}

	db := &db.SQLWrapper{
		DB: bun.NewDB(s, sqlitedialect.New()),
	}

	// Create a new seaport contract listiner
	sl, err := listener.New()
	if err != nil {
		log.Fatalf("Failed to create new SeaportListener: %v", err.Error())
		return err
	}

	// Start the seaport listener
	sl.Start(wg, db, &n.Host)

	// Create a new DHT
	dht := kad.NewDHT(context.Background(), n.Host, nil)

	// Create a new pubsub
	ps, err := pubsub.NewGossipSub(context.Background(), n.Host)
	if err != nil {
		log.Fatalf("Failed to create gossipsub: %v", err.Error())
		return err
	}

	mt, _ := ps.Join("gossipsub:message")
	sub, err := mt.Subscribe(func(subscription *pubsub.Subscription) error {
		log.Printf("Subscription Data: %v", subscription)

		return nil
	})

	if err != nil {
		log.Fatalf("Failed to subscribe to gossipsub: %v", err.Error())
		return err
	}

	handleSub(wg, sub, db)

	err = dht.Bootstrap(context.Background())
	if err != nil {
		log.Fatalf("Failed to bootstrap DHT: %v", err.Error())
		return err
	}

	return nil
}

// TODO: Handle errors
func handleSub(wg *sync.WaitGroup, sub *pubsub.Subscription, database *db.SQLWrapper) {
	wg.Add(1)

	go func() {
		for {
			msg, err := sub.Next(context.Background())
			if err != nil {
				log.Printf("Failed to get next message: %v", err.Error())
				sub.Cancel()
				break
			}

			log.Printf("Message: %v", msg)

			// Save message to the database
			database.Lock()
			defer database.Unlock()

			_, err = database.DB.NewInsert().Model(msg).Exec(context.Background())
			if err != nil {
				log.Printf("Failed to save message to the database: %v", err.Error())
				break
			}

		}

		wg.Done()
	}()
}
