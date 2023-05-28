package service

import (
	"context"
	"reflect"

	"github.com/thanhfphan/eventstore/domain/repos"
	"github.com/thanhfphan/eventstore/pkg/errors"
	"github.com/thanhfphan/eventstore/pkg/ev"
)

var _ AggregateStore = (*aggregateStore)(nil)

type AggregateStore interface {
	Get(ctx context.Context, aggregateID string, agg ev.Aggregate) error
	Save(ctx context.Context, agg ev.Aggregate) error
}

type aggregateStore struct {
	eventRepo     repos.EventRepo
	aggregateRepo repos.AggregateRepo
}

func NewAggregateStore(repos repos.Repos) AggregateStore {
	return &aggregateStore{
		eventRepo:     repos.Event(),
		aggregateRepo: repos.Aggregate(),
	}
}

// Get fetches the events and build up the aggregate
func (as *aggregateStore) Get(ctx context.Context, aggregateID string, agg ev.Aggregate) error {
	if reflect.ValueOf(agg).Kind() != reflect.Ptr {
		return errors.New("aggregate must to be a pointer")
	}

	aggType := reflect.TypeOf(agg).Elem().Name()
	root := agg.Root()
	root.SetAggregateType(aggType)

	events, err := as.eventRepo.Get(ctx, aggregateID, root)
	if err != nil {
		return err
	}

	if len(events) == 0 {
		return errors.New("not found aggregate has id=%s", aggregateID)
	}

	for _, item := range events {
		root.LoadFromHistory(agg, []ev.Event{item})
	}

	return nil
}

func (as *aggregateStore) Save(ctx context.Context, agg ev.Aggregate) error {
	root := agg.Root()
	aggType := reflect.TypeOf(agg).Elem().Name()
	root.SetAggregateType(aggType)

	err := as.aggregateRepo.CreateIfNotExist(ctx, root.AggregateID(), aggType)
	if err != nil {
		return err
	}

	if !as.aggregateRepo.CheckAndUpdateVersion(ctx, agg) {
		return errors.New("Optimistic concurrency control failed id=%s, expectedVersion=%d, newversion=%d",
			agg.Root().AggregateID(), agg.Root().BaseVersion(), agg.Root().Version())
	}

	events := root.Events()
	for _, event := range events {
		err := as.eventRepo.Append(ctx, event)
		if err != nil {
			return err
		}
	}

	root.Update()

	return nil
}
