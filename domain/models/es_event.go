package models

import "github.com/thanhfphan/eventstore/pkg/dtype"

type Event struct {
	ID            int64      `json:"id"`
	TransactionID string     `json:"transaction_id"` // xid8
	AggregateID   string     `json:"aggregate_id"`
	Version       int        `json:"version"`
	Type          string     `json:"type"`
	Data          dtype.JSON `json:"data"`
}
