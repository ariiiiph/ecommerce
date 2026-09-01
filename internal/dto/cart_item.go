package dto

import "time"

type CreateCartItemRequest struct {
	VariantID int64 `json:"variant_id"`
	Quantity  int   `json:"quantity"`
}

type CartItemResponse struct {
	ID        int64     `json:"id"`
	CartID    int64     `json:"cart_id"`
	VariantID int64     `json:"variant_id"`
	Quantity  int       `json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UpdateCartItemRequest struct {
	Quantity int `json:"quantity"`
}
