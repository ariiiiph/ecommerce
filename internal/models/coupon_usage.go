package models

import "time"

type CouponUsage struct {
	ID int64

	CouponID int64
	UserID   int64
	OrderID  int64

	UsedAt time.Time
}
