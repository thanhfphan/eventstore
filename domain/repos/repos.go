package repos

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/thanhfphan/eventstore/pkg/ev"
)

var _ Repos = (*repos)(nil)

type Repos interface {
	Event() EventRepo
	Aggregate() AggregateRepo
}

type repos struct {
	pool      *pgxpool.Pool
	serialize ev.Serializer
}

func New(pool *pgxpool.Pool, s ev.Serializer) Repos {
	return &repos{
		pool:      pool,
		serialize: s,
	}
}

func (r *repos) Event() EventRepo {
	return NewEvent(r.pool, r.serialize)
}

func (r *repos) Aggregate() AggregateRepo {
	return NewAggregate(r.pool, r.serialize)
}
