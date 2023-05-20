package events

import "time"

type OrderCompleted struct {
	CompletedDate *time.Time `json:"completed_date"`
}
