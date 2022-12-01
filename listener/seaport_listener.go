package listener

import (
	"context"
	"goport/abi"
	"goport/config"
	ms "goport/db"
	"log"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/libp2p/go-libp2p/core/host"
)

type SeaportListener struct {
	Seaport             *abi.Seaport
	WatchCountInc       chan *abi.SeaportCounterIncremented
	WatchOrderCancelled chan *abi.SeaportOrderCancelled
	WatchOrderValidated chan *abi.SeaportOrderValidated
	WatchOrderFulfilled chan *abi.SeaportOrderFulfilled
}

// Creates a new SeaportListener
func New() (*SeaportListener, error) {
	c, err := ethclient.Dial(config.RPC_URL)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err.Error())
		return nil, err
	}

	s, err := abi.NewSeaport(common.HexToAddress("0x00000000006c3852cbEf3e08E8dF289169EdE581"), c)
	if err != nil {
		log.Printf("Failed to create new Seaport: %v", err.Error())
		return nil, err
	}

	return &SeaportListener{
		Seaport:       s,
		WatchCountInc: make(chan *abi.SeaportCounterIncremented),
	}, nil
}

func (sl *SeaportListener) Start(db *ms.SQLWrapper, lp *host.Host) {
	sl.watchCounterIncremented(db)
	sl.watchOrderCancelled(db)
	sl.watchOrderValidated(db)
	sl.watchOrderFulfilled(db)
}

// Use context and mutex
func (sl *SeaportListener) watchCounterIncremented(db *ms.SQLWrapper) {
	go func() {
		sub, err := sl.Seaport.WatchCounterIncremented(&bind.WatchOpts{Context: context.Background()}, sl.WatchCountInc, []common.Address{})
		if err != nil {
			log.Printf("Failed to watch CounterIncremented: %v", err.Error())
			return
		}

	countEventLoop:
		for {
			if err := <-sub.Err(); err != nil {
				log.Printf("Failed to watch CounterIncremented: %v", err.Error())
				break countEventLoop
			}

			e := <-sl.WatchCountInc
			err = db.WriteCounterIncremented(e)
			if err != nil {
				log.Printf("Failed to write CounterIncremented: %v", err.Error())
				return
			}

			log.Printf("CounterIncremented: %v", e)
		}

	}()
}

func (sl *SeaportListener) watchOrderCancelled(db *ms.SQLWrapper) {
	go func() {
		sub, err := sl.Seaport.WatchOrderCancelled(&bind.WatchOpts{Context: context.Background()}, sl.WatchOrderCancelled, []common.Address{}, []common.Address{})
		if err != nil {
			log.Printf("Failed to watch OrderCancelled: %v", err.Error())
			return
		}

	cancelledEventLoop:
		for {
			if err := <-sub.Err(); err != nil {
				log.Printf("Failed to watch OrderCancelled: %v", err.Error())
				break cancelledEventLoop
			}

			e := <-sl.WatchOrderCancelled
			err := db.WriteOrderCancelled(e)
			if err != nil {
				log.Printf("Failed to write OrderCancelled: %v", err.Error())
				return
			}

			log.Printf("OrderCancelled: %v", e)
		}
	}()
}

func (sl *SeaportListener) watchOrderValidated(db *ms.SQLWrapper) {
	go func() {
		sub, err := sl.Seaport.WatchOrderValidated(&bind.WatchOpts{Context: context.Background()}, sl.WatchOrderValidated, []common.Address{}, []common.Address{})
		if err != nil {
			log.Printf("Failed to watch OrderValidated: %v", err.Error())
			return
		}

	validatedEventLoop:
		for {
			if err := <-sub.Err(); err != nil {
				log.Printf("Failed to watch OrderValidated: %v", err.Error())
				break validatedEventLoop
			}

			e := <-sl.WatchOrderValidated
			err := db.WriteOrderValidated(e)
			if err != nil {
				log.Printf("Failed to write OrderValidated: %v", err.Error())
				return
			}

			log.Printf("OrderValidated: %v", e)
		}
	}()
}

func (sl *SeaportListener) watchOrderFulfilled(db *ms.SQLWrapper) {
	go func() {
		sub, err := sl.Seaport.WatchOrderFulfilled(&bind.WatchOpts{Context: context.Background()}, sl.WatchOrderFulfilled, []common.Address{}, []common.Address{})
		if err != nil {
			log.Printf("Failed to watch OrderFulfilled: %v", err.Error())
			return
		}

	fulfilledEventLoop:
		for {
			if err := <-sub.Err(); err != nil {
				log.Printf("Failed to watch OrderFulfilled: %v", err.Error())
				break fulfilledEventLoop
			}

			e := <-sl.WatchOrderFulfilled
			err := db.WriteOrderFulfilled(e)
			if err != nil {
				log.Printf("Failed to write OrderFulfilled: %v", err.Error())
				return
			}

			log.Printf("OrderFulfilled: %v", e)
		}
	}()
}
