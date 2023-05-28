package ev

import (
	"encoding/json"
	"reflect"

	"github.com/thanhfphan/eventstore/pkg/errors"
)

type eventFunc = func() interface{}
type RegisterEventsFunc = func(events ...interface{}) error

var _ Serializer = (*serializer)(nil)

type Serializer interface {
	ToEventsFunc(events ...interface{}) []eventFunc
	RegisterAggregate(agg BaseAggregate) error
	Register(agg Aggregate, eventsFunc []eventFunc) error
	RegisterTypes(agg Aggregate, eventsFunc ...eventFunc) error
	Type(typ, name string) (eventFunc, bool)
	Marshal(v interface{}) ([]byte, error)
	Unmarshal(data []byte, v interface{}) error
}

type serializer struct {
	eventRegister map[string]eventFunc
}

func NewSerializer() Serializer {
	return &serializer{
		eventRegister: make(map[string]func() interface{}),
	}
}

func (s *serializer) ToEventsFunc(events ...interface{}) []eventFunc {
	res := []eventFunc{}
	for _, item := range events {
		res = append(res, eventToFunc(item))
	}

	return res
}

func (s *serializer) RegisterAggregate(agg BaseAggregate) error {
	typ := reflect.TypeOf(agg).Elem().Name()
	if typ == "" {
		return errors.New("not found aggregate")
	}

	fu := func(events ...interface{}) error {
		listF := s.ToEventsFunc(events...)
		for _, f := range listF {
			event := f()
			eName := reflect.TypeOf(event).Elem().Name()
			if eName == "" {
				return errors.New("name of event is missing")
			}

			s.eventRegister[typ+"_"+eName] = f
		}

		return nil
	}

	return agg.RegisterEvents(fu)
}

func (s *serializer) Register(agg Aggregate, eventsFunc []eventFunc) error {
	typ := reflect.TypeOf(agg).Elem().Name()
	if typ == "" {
		return errors.New("not found aggregate")
	}
	if len(eventsFunc) == 0 {
		errors.New("not found any register event")
	}

	for _, f := range eventsFunc {
		event := f()
		eName := reflect.TypeOf(event).Elem().Name()
		if eName == "" {
			return errors.New("event name is missing")
		}
		s.eventRegister[typ+"_"+eName] = f
	}

	return nil
}

func (s *serializer) RegisterTypes(agg Aggregate, eventsFunc ...eventFunc) error {
	return s.Register(agg, eventsFunc)
}

func (s *serializer) Type(typ, name string) (eventFunc, bool) {
	f, ok := s.eventRegister[typ+"_"+name]
	return f, ok
}

func (s *serializer) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (s *serializer) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func eventToFunc(e interface{}) eventFunc {
	return func() interface{} {
		return e
	}
}
