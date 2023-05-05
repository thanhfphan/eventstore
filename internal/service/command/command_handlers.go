package command

import (
	"context"

	"github.com/thanhfphan/eventstore/internal/domain/aggregate"
	"github.com/thanhfphan/eventstore/internal/domain/command"
	"github.com/thanhfphan/eventstore/pkg/logging"
)

var (
	_ CommandHandlers = (*cmdHandlers)(nil)
)

type CommandHandlers interface {
	Handle(ctx context.Context, agg aggregate.AggregateRoot, action command.Command) error
}

type cmdHandlers struct {
}

func NewCommandHandlers() CommandHandlers {
	return &cmdHandlers{}
}

func (c *cmdHandlers) Handle(ctx context.Context, agg aggregate.AggregateRoot, action command.Command) error {
	log := logging.FromContext(ctx)
	log.Debugf("cmdHandlers handle agg=%v, action=%v", agg, action)

	switch cmd := action.(type) {
	case *command.PlaceOrderCmd:
		// sometime we need to add more logic here, that why we use switch stament
		return agg.Process(ctx, cmd)
	default:
		return agg.Process(ctx, cmd)
	}
}
