package ev

import (
	"reflect"
	"time"

	"github.com/google/uuid"
	"github.com/thanhfphan/eventstore/pkg/errors"
)

type AggregateRoot struct {
	aggregateID   string
	version       int
	globalVersion int
	events        []Event
}

func (ar *AggregateRoot) SetID(id string) error {
	if ar.aggregateID != "" {
		return errors.New("aggregate already exists")
	}

	ar.aggregateID = id
	return nil
}

func (ar *AggregateRoot) AggregateID() string {
	return ar.aggregateID
}

func (ar *AggregateRoot) Root() *AggregateRoot {
	return ar
}

// version is the version of aggregate not stored
func (ar *AggregateRoot) Version() int {
	if len(ar.events) > 0 {
		return ar.events[len(ar.events)-1].Version
	}

	return ar.version
}

// globalVersion is the version of aggregate has stored
func (ar *AggregateRoot) GlobalVersion() int {
	return ar.globalVersion
}

func (ar *AggregateRoot) Events() []Event {
	return ar.events
}

// CloneEvents return new slice of events
func (ar *AggregateRoot) CloneEvents() []Event {
	evs := make([]Event, len(ar.events))
	copy(evs, ar.events)
	return evs
}

func (ar *AggregateRoot) IsUnsaved() bool {
	return len(ar.events) > 0
}

// ApplyChange apply data change on aggregate
func (ar *AggregateRoot) ApplyChange(agg Aggregate, data interface{}) {
	ar.ApplyChangeWithMetadata(agg, data, nil)
}

func (ar *AggregateRoot) ApplyChangeWithMetadata(agg Aggregate, data interface{}, metadata map[string]interface{}) {
	if ar.aggregateID == "" {
		ar.aggregateID = uuid.NewString()
	}

	name := reflect.TypeOf(agg).Elem().Name()
	event := Event{
		AggregateID:   ar.aggregateID,
		Version:       ar.nextVersion(),
		AggregateType: name,
		CreatedAt:     time.Now().UTC(),
		Data:          data,
		Metadata:      metadata,
	}

	ar.events = append(ar.events, event)
	agg.Transition(event)
}

// LoadFromHistory build aggregate from list event
func (ar *AggregateRoot) LoadFromHistory(agg Aggregate, events []Event) {
	for _, e := range events {
		agg.Transition(e)
		ar.aggregateID = e.AggregateID
		ar.version = e.Version
		ar.globalVersion = e.GlobalVersion
	}
}

// update update version
func (ar *AggregateRoot) Update() {
	if len(ar.events) > 0 {
		lastEvent := ar.events[len(ar.events)-1]
		ar.version = lastEvent.Version
		ar.globalVersion = lastEvent.GlobalVersion
	}
}

// setInternal set common data to AggregateRoot
func (ar *AggregateRoot) setInternal(id string, version, globalVersion int) {
	ar.aggregateID = id
	ar.version = version
	ar.globalVersion = globalVersion
	ar.events = []Event{}
}

func (ar *AggregateRoot) nextVersion() int {
	return ar.version + 1
}

func (ar *AggregateRoot) path() string {
	return reflect.TypeOf(ar).Elem().PkgPath()
}
