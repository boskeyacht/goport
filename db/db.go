package db

import (
	"context"
	"goport/abi"
	"log"
	"sync"

	"github.com/uptrace/bun"
)

type SQLWrapper struct {
	DB *bun.DB
	sync.Mutex
}

func (s *SQLWrapper) WriteCounterIncremented(event *abi.SeaportCounterIncremented) error {
	ic := &CounterIncremented{
		Counter: event.NewCounter,
		Offerer: event.Offerer,
	}

	_, err := s.DB.NewInsert().Model(ic).Exec(context.Background(), s.DB)
	if err != nil {
		return err
	}

	return nil
}

func (s *SQLWrapper) WriteOrderFulfilled(event *abi.SeaportOrderFulfilled) error {
	f := &FulfilledOrder{
		Hash:          event.OrderHash,
		Offerer:       event.Offerer,
		Zone:          event.Zone,
		Recipient:     event.Recipient,
		Offer:         event.Offer,
		Consideration: event.Consideration,
	}

	_, err := s.DB.NewInsert().Model(f).Exec(context.Background(), s.DB)
	if err != nil {
		log.Printf("Failed to write fulfilled order to database: %v", err.Error())
		return err
	}

	return nil
}

func (s *SQLWrapper) WriteOrderCancelled(event *abi.SeaportOrderCancelled) error {
	o := &CancelledOrder{
		Hash:    event.OrderHash,
		Offerer: event.Offerer,
		Zone:    event.Zone,
	}

	_, err := s.DB.NewInsert().Model(o).Exec(context.Background(), s.DB)
	if err != nil {
		log.Printf("Failed to write cancelled order to database: %v", err.Error())
		return err
	}

	return nil
}

func (s *SQLWrapper) WriteOrderValidated(event *abi.SeaportOrderValidated) error {
	v := &ValidatedOrder{
		Hash:    event.OrderHash,
		Offerer: event.Offerer,
		Zone:    event.Zone,
	}

	_, err := s.DB.NewInsert().Model(v).Exec(context.Background(), s.DB)
	if err != nil {
		log.Printf("Failed to write validated order to database: %v", err.Error())
		return err
	}

	return nil
}
