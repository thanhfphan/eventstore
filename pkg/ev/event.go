package ev

import (
	"time"

	"github.com/thanhfphan/eventstore/pkg/dtype"
)

type Event struct {
	ID          int64       `json:"id"`
	AggregateID string      `json:"aggregate_id"`
	Version     int         `json:"version"`
	EventType   string      `json:"event_type"`
	Data        interface{} `json:"data"`
	Metadata    dtype.JSON  `json:"metadata"`

	CreatedAt time.Time `json:"created_at"`
}
