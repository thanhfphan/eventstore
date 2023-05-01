package models

import "github.com/thanhfphan/eventstore/pkg/dtype"

type EsAggregateSnapshot struct {
	AggregateID string     `json:"aggregate_id"` // UUID
	Version     int        `json:"version"`
	Data        dtype.JSON `json:"data"`
}
