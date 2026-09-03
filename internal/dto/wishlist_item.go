package dto

import "time"

type CreateWishlistItemRequest struct {
	ProductID int64 `json:"product_id"`
}

type WishlistItemResponse struct {
	ID         int64     `json:"id"`
	WishlistID int64     `json:"wishlist_id"`
	ProductID  int64     `json:"product_id"`
	CreatedAt  time.Time `json:"created_at"`
}
