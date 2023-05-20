package ev

type BaseAggregate interface {
	RegisterEvents(RegisterEventsFunc) error
}

// All aggregate must implement this interface
type Aggregate interface {
	BaseAggregate
	Root() *AggregateRoot
	Transition(e Event)
}
