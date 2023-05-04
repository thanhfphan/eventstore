package service

import (
	"context"

	"github.com/thanhfphan/eventstore/internal/domain/aggregate"
	"github.com/thanhfphan/eventstore/internal/domain/event"
	"github.com/thanhfphan/eventstore/pkg/dtype"
	"github.com/thanhfphan/eventstore/pkg/logging"
)

var (
	_ AggregateStore = (*aggStore)(nil)
)

type AggregateStore interface {
	Save(ctx context.Context, agg *aggregate.Aggregate) ([]*event.Event, error)
	// CreateAggregateSnapshot(ctx context.Context ...) ...
	Read(ctx context.Context, aagT dtype.AggregateType, aggID string, version int) *aggregate.Aggregate
}

type aggStore struct {
}

func NewAggregateStore() AggregateStore {
	return &aggStore{}
}

func (a *aggStore) Save(ctx context.Context, agg *aggregate.Aggregate) ([]*event.Event, error) {
	log := logging.FromContext(ctx)
	log.Infof("starting Save aggregateStore: %v", agg)

	return nil, nil
}

func (a *aggStore) Read(ctx context.Context, aggT dtype.AggregateType, aggID string, version int) *aggregate.Aggregate {
	log := logging.FromContext(ctx)
	log.Infof("starting READ aggregate with type=%v, ID=%v, version=%d", aggT, aggID, version)

	// TODO: implement load from snapshot

	return nil
}
