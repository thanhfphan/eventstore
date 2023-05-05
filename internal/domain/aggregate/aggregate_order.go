package aggregate

import (
	"context"
	"time"

	"github.com/thanhfphan/eventstore/internal/domain/command"
	"github.com/thanhfphan/eventstore/pkg/errors"
)

var _ AggregateRoot = (*OrderAggregate)(nil)

type OrderAggregate struct {
	Aggregate

	Version       int        `json:"version"`
	Status        string     `json:"status"`
	CustomerID    int64      `json:"customer_id"`
	Price         float64    `json:"price"`
	PlacedDate    *time.Time `json:"place_date"`
	AcceptedDate  *time.Time `json:"accepted_date"`
	CancelledDate *time.Time `json:"cancelled_date"`
	CompletedDate *time.Time `json:"completed_date"`
}

func (a *OrderAggregate) Process(ctx context.Context, c command.Command) error {

	switch cmd := c.(type) {
	case *command.PlaceOrderCmd:
		a.ProcessPlaceCmd(ctx, cmd)
	default:
		return errors.New("OrderAggregate cant handle command=%v", c)
	}

	return nil
}

func (a *OrderAggregate) ProcessPlaceCmd(ctx context.Context, cmd *command.PlaceOrderCmd) {

}
