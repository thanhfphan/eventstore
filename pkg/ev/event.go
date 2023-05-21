package ev

import (
	"reflect"
	"time"

	"github.com/thanhfphan/eventstore/pkg/dtype"
)

type Event struct {
	AggregateID   string      `json:"aggregate_id"`
	Version       int         `json:"version"`
	GlobalVersion int         `json:"global_version"` // consider remove this one
	AggregateType string      `json:"aggregate_type"`
	Data          interface{} `json:"data"`
	Metadata      dtype.JSON  `json:"metadata"`

	CreatedAt time.Time `json:"created_at"`
}

func (e Event) Name() string {
	if e.Data == nil {
		return ""
	}

	return reflect.TypeOf(e.Data).Elem().Name()
}
