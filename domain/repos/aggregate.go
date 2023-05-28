package repos

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/thanhfphan/eventstore/pkg/ev"
	"github.com/thanhfphan/eventstore/pkg/logging"
)

var (
	_ AggregateRepo = (*aggregateRepo)(nil)
)

type AggregateRepo interface {
	CreateIfNotExist(ctx context.Context, id, typ string) error
	CheckAndUpdateVersion(ctx context.Context, agg ev.Aggregate) bool
}

type aggregateRepo struct {
	pool *pgxpool.Pool
}

func NewAggregate(pool *pgxpool.Pool) AggregateRepo {
	return &aggregateRepo{
		pool: pool,
	}
}

func (r *aggregateRepo) CreateIfNotExist(ctx context.Context, id, typ string) error {
	log := logging.FromContext(ctx)
	log.Debugf("CreateIfNotExist id=%s, type=%s", id, typ)

	result, err := r.pool.Exec(ctx, `
		INSERT INTO es_aggregate(id, version, aggregate_type)	
		VALUES($1, 0, $2)
		ON CONFLICT DO NOTHING
	`, id, typ)
	if err != nil {
		log.Warnf("ExecCreateIfNotExist failed with err=%v", err)
		return err
	}

	log.Debugf("CreateIfNotExist got rowsAffected=%d", result.RowsAffected())

	return nil
}

func (r *aggregateRepo) CheckAndUpdateVersion(ctx context.Context, agg ev.Aggregate) bool {
	log := logging.FromContext(ctx)
	log.Debugf("CheckAndUpdateVersion agg=%+v", agg)

	root := agg.Root()

	aggregateId := root.AggregateID()
	expectedVersion := root.BaseVersion()
	newVersion := root.Version()

	result, err := r.pool.Exec(ctx, `
		UPDATE es_aggregate
		SET version = $1
		WHERE id = $2
			AND version = $3
	`, newVersion, aggregateId, expectedVersion)
	if err != nil {
		log.Warnf("check aggregate version failed with err=%v", err)
		return false
	}

	return result.RowsAffected() > 0
}
