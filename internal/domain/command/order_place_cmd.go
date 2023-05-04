package command

import (
	"github.com/google/uuid"
	"github.com/thanhfphan/eventstore/pkg/dtype"
)

type PlaceOrderCmd struct {
	aggType dtype.AggregateType
	aggID   string

	CustomerID int64
	Price      float64
}

func NewPlaceOrderCmd(customerID int64, price float64) Command {
	return &PlaceOrderCmd{
		aggType:    dtype.AggregateTypeOrder,
		aggID:      uuid.NewString(),
		CustomerID: customerID,
		Price:      price,
	}
}

func (c *PlaceOrderCmd) AggregateID() string {
	return c.aggID
}

func (c *PlaceOrderCmd) AggregateType() dtype.AggregateType {
	return c.aggType
}
