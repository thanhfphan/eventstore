package aggregates

import (
	"time"

	"github.com/thanhfphan/eventstore/domain/events"
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
		events.OrderPlaced{},
		events.OrderAccepted{},
		events.OrderCompleted{},
		events.OrderCancelled{},
	)
}

func (o *OrderAggregate) Transition(event ev.Event) {
	// FIXME: cant cast for now
	switch e := event.Data.(type) {
	case *events.OrderPlaced:
		o.OrderID = e.OrderID
		o.CustomerID = e.CustomerID
		o.Price = e.Price
		o.PlacedDate = e.PlacedDate
		o.Status = "placed"
	case *events.OrderAccepted:
		o.AcceptedDate = e.AcceptedDate
		o.Status = "accepted"
	case *events.OrderCompleted:
		o.CompletedDate = e.CompletedDate
		o.Status = "completed"
	case *events.OrderCancelled:
		o.CancelledDate = e.CancelledDate
		o.Status = "cancelled"
	}
}
func CreateOrderAggregate(customerID int64, price float64, date time.Time) *OrderAggregate {
	agg := OrderAggregate{}
	agg.ApplyChange(&agg, &events.OrderPlaced{
		CustomerID: customerID,
		Price:      price,
		PlacedDate: date,
	})

	return &agg
}

func (o *OrderAggregate) RecordAccepted(date *time.Time) {
	o.ApplyChange(o, &events.OrderAccepted{
		AcceptedDate: date,
	})
}

func (o *OrderAggregate) RecordCompleted(date *time.Time) {
	o.ApplyChange(o, &events.OrderCompleted{
		CompletedDate: date,
	})
}

func (o *OrderAggregate) RecordCancelled(date time.Time) {
	o.ApplyChange(o, &events.OrderCancelled{
		CancelledDate: &date,
	})
}
