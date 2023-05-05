package aggregate

import (
	"context"

	"github.com/thanhfphan/eventstore/internal/domain/command"
	"github.com/thanhfphan/eventstore/internal/domain/event"
)

type AggregateRoot interface {
	Process(ctx context.Context, cmd command.Command) error
}

type Aggregate struct {
	ID      string `json:"id"`
	Version int    `json:"version"`
	Type    string `json:"type"`

	Changes     []*event.Event `json:"-"`
	BaseVersion int            `json:"-"`
}
