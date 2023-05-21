package repos

import (
	"context"
	"database/sql"

	"github.com/thanhfphan/eventstore/pkg/logging"
)

var (
	_ AggregateRepo = (*aggregateRepo)(nil)
)

type AggregateRepo interface {
	CreateIfNotExist(ctx context.Context, id, typ string) error
}

type aggregateRepo struct {
	db *sql.DB
}

func NewAggregate(db *sql.DB) AggregateRepo {
	return &aggregateRepo{
		db: db,
	}
}

func (r *aggregateRepo) CreateIfNotExist(ctx context.Context, id, typ string) error {
	log := logging.FromContext(ctx)
	log.Debugf("CreateIfNotExist id=%s, type=%s", id, typ)

	result, err := r.db.ExecContext(ctx, `
		INSERT INTO es_aggregate(id, version, type)	
		VALUES($1, 0, $2)
		ON CONFLICT DO NOTHING
	`, id, typ)
	if err != nil {
		log.Warnf("ExecCreateIfNotExist failed with err=%v", err)
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		log.Warnf("create if not exist failed with err=%v", err)
		return err
	}

	_ = rows

	return nil
}
