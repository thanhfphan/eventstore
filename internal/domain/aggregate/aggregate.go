package aggregate

import (
	"context"

	"github.com/thanhfphan/eventstore/internal/domain/command"
	"github.com/thanhfphan/eventstore/internal/domain/event"
	"github.com/thanhfphan/eventstore/pkg/errors"
)

type Aggregate struct {
	ID      string `json:"id"`
	Version int    `json:"version"`
	Type    string `json:"type"`

	Changes     []*event.Event `json:"-"`
	BaseVersion int            `json:"-"`
}

func (a *Aggregate) Process(ctx context.Context, cmd command.Command) error {
	f := getInstance(cmd.AggregateType())
	if f != nil {
		return f.Process(ctx, a, cmd)
	}

	return errors.New("do not found func handle AggregateType=%v", cmd.AggregateType())
}
