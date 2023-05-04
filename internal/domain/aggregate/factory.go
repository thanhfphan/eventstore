package aggregate

import "github.com/thanhfphan/eventstore/pkg/dtype"

func getInstance(t dtype.AggregateType) AggregateBase {
	switch t {
	case dtype.AggregateTypeOrder:
		return &OrderAggregate{}
	default:
		return nil
	}
}
