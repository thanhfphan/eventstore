package repos

import (
	"context"
	"database/sql"

	"github.com/thanhfphan/eventstore/domain/models"
)

var (
	_ EventRepo = (*eventRepo)(nil)
)

type EventRepo interface {
	AppendEvent(ctx context.Context, e *models.Event) (int64, error)
	ReadEvents(ctx context.Context, aggregateID string, fromVersion, toVersion int) ([]*models.Event, error)
}

type eventRepo struct {
	db *sql.DB
}

func NewEvent(db *sql.DB) EventRepo {
	return &eventRepo{
		db: db,
	}
}

func (r *eventRepo) AppendEvent(ctx context.Context, e *models.Event) (int64, error) {
	result, err := r.db.Exec(`
		INSERT INTO es_event(transaction_id, aggregate_id, version, type, data)
		VALUES(pg_current_xact_id(), ?, ?, ?, ?)`, e.AggregateID, e.Version, e.Type, e.Data)
	if err != nil {
		return -1, err
	}

	return result.LastInsertId()
}
func (r *eventRepo) ReadAllEvents(ctx context.Context, aggregateID string) ([]*models.Event, error) {

	rows, err := r.db.QueryContext(ctx, `
			SELECT id, transaction_id::text, type, data
			FROM es_event
			WHERE aggregate_id = ?
			ORDER BY version ASC`, aggregateID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []*models.Event
	for rows.Next() {
		var evt *models.Event
		if err := rows.Scan(&evt.ID, &evt.TransactionID, &evt.Type, &evt.Data); err != nil {
			return nil, err
		}
		records = append(records, evt)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return records, nil
}

func (r *eventRepo) ReadEvents(ctx context.Context, aggregateID string, fromVersion, toVersion int) ([]*models.Event, error) {

	rows, err := r.db.QueryContext(ctx, `
			SELECT id, transaction_id::text, type, data
			FROM es_event
			WHERE aggregate_id = ?
				AND version > ?
				AND version <= ?
			ORDER BY version ASC`, aggregateID, fromVersion, toVersion)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []*models.Event
	for rows.Next() {
		var evt *models.Event
		if err := rows.Scan(&evt.ID, &evt.TransactionID, &evt.Type, &evt.Data); err != nil {
			return nil, err
		}
		records = append(records, evt)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return records, nil
}
