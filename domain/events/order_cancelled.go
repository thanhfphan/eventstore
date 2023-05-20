package events

import "time"

type OrderCancelled struct {
	CancelledDate *time.Time `json:"cancelled_date"`
}
