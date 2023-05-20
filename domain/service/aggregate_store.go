package service

var _ AggregateStore = (*aggregateStore)(nil)

type AggregateStore interface {
}

type aggregateStore struct {
}

func NewAggregateStore() AggregateStore {
	return &aggregateStore{}
}

// TODO
// func Save(

// TODO
// func Get(
