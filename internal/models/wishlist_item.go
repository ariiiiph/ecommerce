package models

import "time"

type WishlistItem struct {
	ID         int64
	WishlistID int64
	ProductID  int64
	CreatedAt  time.Time
}
