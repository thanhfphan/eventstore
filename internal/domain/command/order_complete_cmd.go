package command

import (
	"github.com/thanhfphan/eventstore/pkg/dtype"
)

type CompleteOrderCmd struct {
	aggType dtype.AggregateType
	aggID   string
}

func NewCompleteOrderCmd(aggregateID string) Command {
	return &CompleteOrderCmd{
		aggType: dtype.AggregateTypeOrder,
		aggID:   aggregateID,
	}
}

func (c *CompleteOrderCmd) AggregateID() string {
	return c.aggID
}

func (c *CompleteOrderCmd) AggregateType() dtype.AggregateType {
	return c.aggType
}
