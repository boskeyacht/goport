package listener

import (
	"context"
	"goport/abi"
	"goport/config"
	ms "goport/db"
	"log"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
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

func (sl *SeaportListener) Start(wg *sync.WaitGroup, db *ms.SQLWrapper) {
	sl.watchCounterIncremented(wg, db)
	sl.watchOrderCancelled(wg, db)
	sl.watchOrderValidated(wg, db)
	sl.watchOrderFulfilled(wg, db)
}

// Watches the Seaport contract for a counter incremented event and writes it to the database
func (sl *SeaportListener) watchCounterIncremented(wg *sync.WaitGroup, db *ms.SQLWrapper) {
	wg.Add(1)
	go func() {
		sub, err := sl.Seaport.WatchCounterIncremented(&bind.WatchOpts{Context: context.Background()}, sl.WatchCountInc, []common.Address{})
		if err != nil {
			log.Printf("Failed to watch CounterIncremented: %v", err.Error())
			return
		}

		for {
			if err := <-sub.Err(); err != nil {
				log.Printf("Failed to watch CounterIncremented: %v", err.Error())
				break
			}

			e := <-sl.WatchCountInc

			db.Lock()

			err = db.WriteCounterIncremented(e)
			if err != nil {
				log.Printf("Failed to write CounterIncremented: %v", err.Error())
				break
			}

			db.Unlock()
			log.Printf("CounterIncremented: %v", e.Raw.Address.String())
		}

		wg.Done()
	}()

}

// Watches the Seaport contract for a order cancelled event and writes it to the database
func (sl *SeaportListener) watchOrderCancelled(wg *sync.WaitGroup, db *ms.SQLWrapper) {
	wg.Add(1)
	go func() {
		sub, err := sl.Seaport.WatchOrderCancelled(&bind.WatchOpts{Context: context.Background()}, sl.WatchOrderCancelled, []common.Address{}, []common.Address{})
		if err != nil {
			log.Printf("Failed to watch OrderCancelled: %v", err.Error())
			return
		}

		for {
			if err := <-sub.Err(); err != nil {
				log.Printf("Failed to watch OrderCancelled: %v", err.Error())
				break
			}

			e := <-sl.WatchOrderCancelled

			db.Lock()

			err := db.WriteOrderCancelled(e)
			if err != nil {
				log.Printf("Failed to write OrderCancelled: %v", err.Error())
				break
			}

			db.Unlock()
			log.Printf("OrderCancelled: %v", e.Raw.Address.String())
		}

		wg.Done()
	}()
}

// Watches the Seaport contract for a order validated event and writes it to the database
func (sl *SeaportListener) watchOrderValidated(wg *sync.WaitGroup, db *ms.SQLWrapper) {
	wg.Add(1)
	go func() {
		sub, err := sl.Seaport.WatchOrderValidated(&bind.WatchOpts{Context: context.Background()}, sl.WatchOrderValidated, []common.Address{}, []common.Address{})
		if err != nil {
			log.Printf("Failed to watch OrderValidated: %v", err.Error())
			return
		}

		for {
			if err := <-sub.Err(); err != nil {
				log.Printf("Failed to watch OrderValidated: %v", err.Error())
				break
			}

			e := <-sl.WatchOrderValidated

			db.Lock()

			err := db.WriteOrderValidated(e)
			if err != nil {
				log.Printf("Failed to write OrderValidated: %v", err.Error())
				break
			}

			db.Unlock()
			log.Printf("OrderValidated: %v", e.Raw.Address.String())
		}

		wg.Done()
	}()
}

// Watches the Seaport contract for a order fulfilled event and writes it to the database
func (sl *SeaportListener) watchOrderFulfilled(wg *sync.WaitGroup, db *ms.SQLWrapper) {
	wg.Add(1)

	go func() {
		sub, err := sl.Seaport.WatchOrderFulfilled(&bind.WatchOpts{Context: context.Background()}, sl.WatchOrderFulfilled, []common.Address{}, []common.Address{})
		if err != nil {
			log.Printf("Failed to watch OrderFulfilled: %v", err.Error())
			return
		}

		for {
			if err := <-sub.Err(); err != nil {
				log.Printf("Failed to watch OrderFulfilled: %v", err.Error())
				break
			}

			e := <-sl.WatchOrderFulfilled

			db.Lock()

			err := db.WriteOrderFulfilled(e)
			if err != nil {
				log.Printf("Failed to write OrderFulfilled: %v", err.Error())
				break
			}

			db.Unlock()
			log.Printf("OrderFulfilled: %v", e.Raw.Address.String())

		}

		wg.Done()
	}()
}
