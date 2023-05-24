package models

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Event struct {
	ID            int64             `json:"id"`
	TransactionID string            `json:"transaction_id"` // xid8
	AggregateID   string            `json:"aggregate_id"`
	Version       int               `json:"version"`
	Type          string            `json:"type"`
	Data          pgtype.JSONBCodec `json:"data"`
}
