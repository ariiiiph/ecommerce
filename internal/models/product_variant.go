package models

import "time"

type ProductVariant struct {
	ID            int64     `json:"id"`
	ProductID     int64     `json:"product_id"`
	SKU           string    `json:"sku"`
	Price         float64   `json:"price"`
	DiscountPrice *float64  `json:"discount_price,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
