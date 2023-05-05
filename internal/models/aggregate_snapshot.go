package models

import "github.com/thanhfphan/eventstore/pkg/dtype"

type AggregateSnapshot struct {
	AggregateType string     `json:"aggregate_type"`
	Data          dtype.JSON `json:"data"`
}
