package models

import "time"

type Order struct {
	OrderID       string     `json:"order_id"`
	CustomerID    int64      `json:"customer_id"`
	Price         float64    `json:"price"`
	Status        string     `json:"status"`
	PlacedDate    *time.Time `json:"placed_date"`
	AcceptedDate  *time.Time `json:"accepted_date"`
	CompletedDate *time.Time `json:"completed_date"`
	CancelledDate *time.Time `json:"cancelled_date"`
}
