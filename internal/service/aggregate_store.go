package service

import (
	"context"

	"github.com/thanhfphan/eventstore/internal/domain/aggregate"
	"github.com/thanhfphan/eventstore/internal/domain/event"
	"github.com/thanhfphan/eventstore/internal/repos"
	"github.com/thanhfphan/eventstore/pkg/dtype"
	"github.com/thanhfphan/eventstore/pkg/logging"
)

var (
	_ AggregateStore = (*aggStore)(nil)
)

type AggregateStore interface {
	Save(ctx context.Context, agg aggregate.AggregateRoot) ([]*event.Event, error)
	// CreateAggregateSnapshot(ctx context.Context ...) ...
	Read(ctx context.Context, aagT dtype.AggregateType, aggID string) aggregate.AggregateRoot
}

type aggStore struct {
	eventRepo     repos.EventRepo
	aggregateRepo repos.AggregateRepo
}

func NewAggregateStore(repos repos.Repos) AggregateStore {
	return &aggStore{
		eventRepo:     repos.Event(),
		aggregateRepo: repos.Aggregate(),
	}
}

func (a *aggStore) Save(ctx context.Context, agg aggregate.AggregateRoot) ([]*event.Event, error) {
	log := logging.FromContext(ctx)
	log.Debugf("starting Save aggregateStore: %v", agg)

	return nil, nil
}

func (a *aggStore) Read(ctx context.Context, aggT dtype.AggregateType, aggID string) aggregate.AggregateRoot {
	log := logging.FromContext(ctx)
	log.Debugf("starting READ aggregate with type=%v, ID=%v", aggT, aggID)

	aggregate := a.readFromSnapshot(ctx, aggID, -1)
	if aggregate != nil {
		return aggregate
	}

	log.Debugf("AggregateID=%s not found", aggID)

	aggregate = a.readFromEvents(ctx, aggT, aggID, -1)
	log.Debugf("Read aggreate has data=%v", aggregate)

	return aggregate
}

func (a *aggStore) readFromEvents(ctx context.Context,
	aggT dtype.AggregateType,
	aggID string,
	version int) aggregate.AggregateRoot {
	log := logging.FromContext(ctx)
	log.Debugf("starting READ aggregate with, ID=%v", aggID)

	return nil
}

func (a *aggStore) readFromSnapshot(ctx context.Context, aggID string, version int) aggregate.AggregateRoot {
	log := logging.FromContext(ctx)
	log.Debugf("starting READ aggregate with, ID=%v", aggID)

	agg, err := a.aggregateRepo.ReadAggregateSnapshot(ctx, aggID, version)
	if err != nil {
		return nil
	}

	// TODO: check version snapshot

	return agg
}
