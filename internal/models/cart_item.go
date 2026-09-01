package models

import "time"

type CartItem struct {
	ID        int64
	CartID    int64
	VariantID int64
	Quantity  int
	CreatedAt time.Time
	UpdatedAt time.Time
}
