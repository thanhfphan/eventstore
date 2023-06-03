package repos

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/thanhfphan/eventstore/pkg/errors"
	"github.com/thanhfphan/eventstore/pkg/ev"
	"github.com/thanhfphan/eventstore/pkg/logging"
)

var (
	_ EventRepo = (*eventRepo)(nil)
)

type EventRepo interface {
	Get(ctx context.Context, aggregateID string, fromVer, toVer int, root *ev.AggregateRoot) ([]ev.Event, error)
	Append(context.Context, ev.Event) error
}

type eventRepo struct {
	pool      *pgxpool.Pool
	serialize ev.Serializer
}

func NewEvent(pool *pgxpool.Pool, s ev.Serializer) EventRepo {
	return &eventRepo{
		pool:      pool,
		serialize: s,
	}
}

func (r *eventRepo) Get(ctx context.Context, aggregateID string, fromVer, toVer int, root *ev.AggregateRoot) ([]ev.Event, error) {
	log := logging.FromContext(ctx)
	log.Debugf("Starting GetAfter aggregateID=%s, fromVersion=%d, toVersion=%d", aggregateID, fromVer, toVer)

	rows, err := r.pool.Query(ctx, `
			SELECT id, aggregate_id, event_type, version, data
			FROM es_event
			WHERE aggregate_id = $1
				AND ($2 = 0 OR version > $2)
				AND ($3 = 0 OR version <= $3)
			ORDER BY version ASC`, aggregateID, fromVer, toVer)
	if err != nil {
		log.Warnf("query get after event failed with err=%v", err)
		return nil, err
	}
	defer rows.Close()

	var records []ev.Event
	for rows.Next() {
		var evt ev.Event
		var data string
		if err := rows.Scan(&evt.ID, &evt.AggregateID, &evt.EventType, &evt.Version, &data); err != nil {
			log.Warnf("scan event in getafter failed with err=%v", err)
			return nil, err
		}

		f, ok := r.serialize.Type(root.AggregateType(), evt.EventType)
		if !ok {
			log.Warnf("for some reason cant serialize event with type: %s_%s", root.AggregateType(), evt.EventType)
			continue
		}

		eventData := f()
		err := r.serialize.Unmarshal([]byte(data), &eventData)
		if err != nil {
			return nil, errors.New("unmarshal event failed with err=%v", err)
		}

		evt.Data = eventData
		records = append(records, evt)
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

	eData, err := r.serialize.Marshal(e.Data)
	if err != nil {
		return err
	}

	result, err := r.pool.Exec(ctx, `
		INSERT INTO es_event(transaction_id, aggregate_id, event_type, version, data, metadata)
		VALUES(pg_current_xact_id(), $1, $2, $3, $4, $5)`,
		e.AggregateID, e.EventType, e.Version, eData, e.Metadata)
	if err != nil {
		log.Warnf("append event failed with err=%v", err)
		return err
	}

	if !result.Insert() {
		return errors.New("insert failed")
	}

	return nil
}
