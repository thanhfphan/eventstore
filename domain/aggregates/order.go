package aggregates

import (
	"time"

	"github.com/thanhfphan/eventstore/domain/models"
	"github.com/thanhfphan/eventstore/pkg/ev"
)

var _ ev.Aggregate = (*OrderAggregate)(nil)

type OrderAggregate struct {
	ev.AggregateRoot
	models.Order
}

func (o *OrderAggregate) RegisterEvents(f ev.RegisterEventsFunc) error {
	return f(
		&OrderPlaced{},
		&OrderAccepted{},
		&OrderCompleted{},
		&OrderCancelled{},
	)
}

func (o *OrderAggregate) Transition(event ev.Event) {
	// FIXME: cant cast for now
	switch e := event.Data.(type) {
	case *OrderPlaced:
		o.OrderID = e.OrderID
		o.CustomerID = e.CustomerID
		o.Price = e.Price
		o.PlacedDate = e.PlacedDate
		o.Status = "placed"
	case *OrderAccepted:
		o.AcceptedDate = e.AcceptedDate
		o.Status = "accepted"
	case *OrderCompleted:
		o.CompletedDate = e.CompletedDate
		o.Status = "completed"
	case *OrderCancelled:
		o.CancelledDate = e.CancelledDate
		o.Status = "cancelled"
	}
}
func CreateOrderAggregate(customerID int64, price float64, date time.Time) *OrderAggregate {
	agg := OrderAggregate{}
	agg.ApplyChange(&agg, &OrderPlaced{
		CustomerID: customerID,
		Price:      price,
		PlacedDate: date,
	})

	return &agg
}

func (o *OrderAggregate) RecordAccepted(date *time.Time) {
	o.ApplyChange(o, &OrderAccepted{
		AcceptedDate: date,
	})
}

func (o *OrderAggregate) RecordCompleted(date *time.Time) {
	o.ApplyChange(o, &OrderCompleted{
		CompletedDate: date,
	})
}

func (o *OrderAggregate) RecordCancelled(date time.Time) {
	o.ApplyChange(o, &OrderCancelled{
		CancelledDate: &date,
	})
}
