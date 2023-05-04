package aggregate

import (
	"context"

	"github.com/thanhfphan/eventstore/internal/domain/command"
)

type AggregateBase interface {
	Process(ctx context.Context, agg *Aggregate, cmd command.Command) error
}
