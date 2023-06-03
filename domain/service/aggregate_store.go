package service

import (
	"context"
	"reflect"

	"github.com/thanhfphan/eventstore/domain/repos"
	"github.com/thanhfphan/eventstore/pkg/errors"
	"github.com/thanhfphan/eventstore/pkg/ev"
	"github.com/thanhfphan/eventstore/pkg/logging"
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
	log := logging.FromContext(ctx)

	if reflect.ValueOf(agg).Kind() != reflect.Ptr {
		return errors.New("aggregate must to be a pointer")
	}

	aggType := reflect.TypeOf(agg).Elem().Name()
	agg.Root().SetAggregateType(aggType)

	// FIXME: get snapshot configurable
	has := as.getFromSnapshot(ctx, aggregateID, agg)
	if !has {
		log.Debugf("not found in snapshot")
		return as.getFromEvents(ctx, aggregateID, agg)
	}

	return nil
}

func (as *aggregateStore) Save(ctx context.Context, agg ev.Aggregate) error {
	log := logging.FromContext(ctx)

	root := agg.Root()
	aggType := reflect.TypeOf(agg).Elem().Name()
	root.SetAggregateType(aggType)

	err := as.aggregateRepo.CreateIfNotExist(ctx, root.AggregateID(), aggType)
	if err != nil {
		return err
	}

	if !as.aggregateRepo.CheckAndUpdateVersion(ctx, agg) {
		return errors.New("Optimistic concurrency control failed id=%s, expectedVersion=%d, newversion=%d",
			root.AggregateID(), root.BaseVersion(), root.Version())
	}

	events := root.Events()
	for _, event := range events {
		err := as.eventRepo.Append(ctx, event)
		if err != nil {
			return err
		}

		// FIXME: configurable
		err = as.createSnapshot(ctx, agg)
		if err != nil {
			log.Warnf("create snapshot aggId=%s, aggType=%s failed err=%v", root.AggregateID(), root.AggregateType(), err)
		}
	}

	root.Update()

	return nil
}

func (as *aggregateStore) getFromSnapshot(ctx context.Context, aggregateID string, agg ev.Aggregate) bool {
	log := logging.FromContext(ctx)
	has := as.aggregateRepo.ReadSnapshot(ctx, aggregateID, agg.Root().Version(), agg)
	if !has {
		return false
	}

	root := agg.Root()
	if root.BaseVersion() < root.Version() {
		events, err := as.eventRepo.Get(ctx, aggregateID, root.BaseVersion(), root.Version(), root)
		if err != nil {
			log.Warnf("get event fromVer=%d to toVer=%s failed err=%v", root.BaseVersion(), root.Version(), err)
			return false
		}

		log.Debugf("event %v", events)
		for _, item := range events {
			root.LoadFromHistory(agg, []ev.Event{item})
		}
	}

	return true
}

func (as *aggregateStore) getFromEvents(ctx context.Context, aggregateID string, agg ev.Aggregate) error {
	root := agg.Root()

	events, err := as.eventRepo.Get(ctx, aggregateID, 0, 0, root)
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

func (as *aggregateStore) createSnapshot(ctx context.Context, agg ev.Aggregate) error {
	log := logging.FromContext(ctx)

	nthEvent := 3 // FIXME: configurable
	root := agg.Root()
	if root.BaseVersion()%nthEvent == 0 {
		log.Infof("create snapshot of aggId=%s, aggType=%s", root.AggregateID(), root.AggregateType())
		return as.aggregateRepo.CreateSnapshot(ctx, agg)
	}

	return nil
}
