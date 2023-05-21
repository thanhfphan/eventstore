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

	root := agg.Root()
	events, err := as.eventRepo.GetAfter(ctx, aggregateID, root.Version())
	if err != nil {
		return err
	}

	for _, item := range events {
		e := ev.Event{
			Version:       item.Version,
			Data:          item.Data,
			AggregateID:   item.AggregateID,
			AggregateType: item.Type,
		}

		root.LoadFromHistory(agg, []ev.Event{e})
	}

	return nil
}

func (as *aggregateStore) Save(ctx context.Context, agg ev.Aggregate) error {
	root := agg.Root()

	aggType := reflect.TypeOf(agg).Elem().Name()

	err := as.aggregateRepo.CreateIfNotExist(ctx, root.AggregateID(), aggType)
	if err != nil {
		return err
	}

	// TODO: check version

	events := root.Events()
	for _, event := range events {
		err := as.eventRepo.Append(ctx, event)
		if err != nil {
			return err
		}

		// events[i].GlobalVersion = newEvent.ID
	}

	root.Update()

	return nil
}
