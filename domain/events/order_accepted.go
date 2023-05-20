package events

import "time"

type OrderAccepted struct {
	AcceptedDate *time.Time `json:"accepted_date"`
}
