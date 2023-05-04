package command

import (
	"github.com/thanhfphan/eventstore/pkg/dtype"
)

type Command interface {
	AggregateID() string
	AggregateType() dtype.AggregateType
}
