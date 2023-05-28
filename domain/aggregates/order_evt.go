package aggregates

import "time"

type OrderPlaced struct {
	OrderID    string    `json:"order_id"`
	CustomerID int64     `json:"customer_id"`
	Price      float64   `json:"price"`
	PlacedDate time.Time `json:"placed_date"`
}

type OrderCancelled struct {
	CancelledDate *time.Time `json:"cancelled_date"`
}

type OrderCompleted struct {
	CompletedDate *time.Time `json:"completed_date"`
}

type OrderAccepted struct {
	AcceptedDate *time.Time `json:"accepted_date"`
}
