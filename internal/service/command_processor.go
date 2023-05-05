package service

import (
	"context"

	"github.com/thanhfphan/eventstore/internal/domain/command"
	"github.com/thanhfphan/eventstore/internal/repos"
	svcCmd "github.com/thanhfphan/eventstore/internal/service/command"
	"github.com/thanhfphan/eventstore/pkg/logging"
)

var (
	_ CommandProcessor = (*cmdProcessor)(nil)
)

type CommandProcessor interface {
	Process(ctx context.Context, cmd command.Command) error
}

type cmdProcessor struct {
	aggStore    AggregateStore
	cmdHandlers svcCmd.CommandHandlers
}

func NewCommandProcessor(repos repos.Repos) CommandProcessor {
	return &cmdProcessor{
		aggStore:    NewAggregateStore(repos),
		cmdHandlers: svcCmd.NewCommandHandlers(),
	}
}

func (c *cmdProcessor) Process(ctx context.Context, cmd command.Command) error {
	log := logging.FromContext(ctx)
	log.Debugf("starting process command: %v", cmd)

	aggregate := c.aggStore.Read(ctx, cmd.AggregateType(), cmd.AggregateID())

	err := c.cmdHandlers.Handle(ctx, aggregate, cmd)
	if err != nil {
		log.Warnf("run cmdHandlers failed with err=%v", err)
		return err
	}

	events, err := c.aggStore.Save(ctx, aggregate)
	if err != nil {
		log.Warnf("save aggStore failed with err=%v", err)
		return err
	}

	// FIXME: might want to publish that event to somewhere else
	_ = events

	return nil
}
