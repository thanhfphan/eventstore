package repos

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ Repos = (*repos)(nil)

type Repos interface {
	Event() EventRepo
	Aggregate() AggregateRepo
}

type repos struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) Repos {
	return &repos{
		pool: pool,
	}
}

func (r *repos) Event() EventRepo {
	return NewEvent(r.pool)
}

func (r *repos) Aggregate() AggregateRepo {
	return NewAggregate(r.pool)
}
