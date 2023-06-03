package repos

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/thanhfphan/eventstore/pkg/errors"
	"github.com/thanhfphan/eventstore/pkg/ev"
	"github.com/thanhfphan/eventstore/pkg/logging"
)

var (
	_ AggregateRepo = (*aggregateRepo)(nil)
)

type AggregateRepo interface {
	CreateIfNotExist(ctx context.Context, id, typ string) error
	CheckAndUpdateVersion(ctx context.Context, agg ev.Aggregate) bool
	CreateSnapshot(ctx context.Context, agg ev.Aggregate) error
	ReadSnapshot(ctx context.Context, aggregateID string, version int, agg ev.Aggregate) bool
}

type aggregateRepo struct {
	pool      *pgxpool.Pool
	serialize ev.Serializer
}

func NewAggregate(pool *pgxpool.Pool, s ev.Serializer) AggregateRepo {
	return &aggregateRepo{
		serialize: s,
		pool:      pool,
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
	log.Debugf("CheckAndUpdateVersion agg=%+v", agg.Root())

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

func (r *aggregateRepo) ReadSnapshot(ctx context.Context, aggregateID string, version int, agg ev.Aggregate) bool {
	log := logging.FromContext(ctx)
	var (
		aggType     string
		data        string
		aggVersion  int
		snapVersion int
	)
	err := r.pool.QueryRow(ctx, `
		SELECT a.aggregate_type, a.version as aggregateVer, eas.version as snapshotVer, eas.data
		FROM es_aggregate_snapshot eas
		JOIN es_aggregate a ON eas.aggregate_id = a.id
		WHERE eas.aggregate_id = $1
			AND eas.version >= $2
		ORDER BY eas.version DESC
		LIMIT 1
		`, aggregateID, version).Scan(&aggType, &aggVersion, &snapVersion, &data)
	if err != nil {
		return false
	}
	root := agg.Root()

	err = r.serialize.Unmarshal([]byte(data), agg)
	if err != nil {
		log.Warnf("serialize aggData failed err=%v", err)
		return false
	}

	root.SetInternal(aggregateID, snapVersion, aggVersion)

	return true
}

func (r *aggregateRepo) CreateSnapshot(ctx context.Context, agg ev.Aggregate) error {
	log := logging.FromContext(ctx)

	root := agg.Root()
	aggregateId := root.AggregateID()
	version := root.Version()
	data, err := r.serialize.Marshal(agg)
	if err != nil {
		return err
	}

	result, err := r.pool.Exec(ctx, `
		INSERT INTO es_aggregate_snapshot (aggregate_id, version, data)
		VALUES($1, $2, $3)`, aggregateId, version, data)
	if err != nil {
		log.Warnf("insert es_aggregate_snapshot failed err=%v", err)
		return err
	}

	if !result.Insert() {
		return errors.New("insert failed")
	}

	return nil
}
