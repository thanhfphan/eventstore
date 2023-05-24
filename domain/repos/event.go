package repos

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/thanhfphan/eventstore/domain/models"
	"github.com/thanhfphan/eventstore/pkg/errors"
	"github.com/thanhfphan/eventstore/pkg/ev"
	"github.com/thanhfphan/eventstore/pkg/logging"
)

var (
	_ EventRepo = (*eventRepo)(nil)
)

type EventRepo interface {
	GetAfter(ctx context.Context, aggregateID string, fromVersion int) ([]*models.Event, error)
	Append(context.Context, ev.Event) error
}

type eventRepo struct {
	pool *pgxpool.Pool
}

func NewEvent(pool *pgxpool.Pool) EventRepo {
	return &eventRepo{
		pool: pool,
	}
}

func (r *eventRepo) GetAfter(ctx context.Context, aggregateID string, fromVersion int) ([]*models.Event, error) {
	log := logging.FromContext(ctx)
	log.Debugf("Starting GetAfter aggregateID=%s, fromVersion=%d", aggregateID, fromVersion)

	rows, err := r.pool.Query(ctx, `
			SELECT id, aggregate_id, type, version, data
			FROM es_event
			WHERE aggregate_id = $1
				AND version > $2
			ORDER BY version ASC`, aggregateID, fromVersion)
	if err != nil {
		log.Warnf("query get after event failed with err=%v", err)
		return nil, err
	}
	defer rows.Close()

	var records []*models.Event
	for rows.Next() {
		var evt models.Event
		if err := rows.Scan(&evt.ID, &evt.AggregateID, &evt.Type, &evt.Version, &evt.Data); err != nil {
			log.Warnf("scan event in getafter failed with err=%v", err)
			return nil, err
		}
		records = append(records, &evt)
	}

	if err := rows.Err(); err != nil {
		log.Warnf("getAfter: rows.err failed", err)
		return nil, err
	}

	return records, nil
}

func (r *eventRepo) Append(ctx context.Context, e ev.Event) error {
	log := logging.FromContext(ctx)
	log.Debugf("Starting append event=%v", e)

	eData, err := json.Marshal(e.Data)
	if err != nil {
		return err
	}

	result, err := r.pool.Exec(ctx, `
		INSERT INTO es_event(transaction_id, aggregate_id, version, type, data)
		VALUES(pg_current_xact_id(), $1, $2, $3, $4)`,
		e.AggregateID, e.Version, e.AggregateType, eData)
	if err != nil {
		log.Warnf("append event failed with err=%v", err)
		return err
	}

	if !result.Insert() {
		return errors.New("insert failed")
	}

	return nil
}
