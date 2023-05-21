package repos

import "database/sql"

var _ Repos = (*repos)(nil)

type Repos interface {
	Event() EventRepo
	Aggregate() AggregateRepo
}

type repos struct {
	db *sql.DB
}

func New(db *sql.DB) Repos {
	return &repos{
		db: db,
	}
}

func (r *repos) Event() EventRepo {
	return NewEvent(r.db)
}

func (r *repos) Aggregate() AggregateRepo {
	return NewAggregate(r.db)
}
