package events

import "time"

type OrderPlaced struct {
	OrderID    string     `json:"order_id"`
	CustomerID int64      `json:"customer_id"`
	Price      float64    `json:"price"`
	PlacedDate *time.Time `json:"placed_date"`
}
