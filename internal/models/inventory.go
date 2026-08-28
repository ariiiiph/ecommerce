package models

import "time"

type Inventory struct {
	VariantID         int64     `json:"variant_id"`
	Quantity          int       `json:"quantity"`
	ReservedQuantity  int       `json:"reserved_quantity"`
	LowStockThreshold int       `json:"low_stock_threshold"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}
