package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type Coupon struct {
	ID                 int64
	Code               string
	DiscountType       string
	DiscountValue      decimal.Decimal
	MinimumOrderAmount *decimal.Decimal
	MaximumDiscount    *decimal.Decimal
	UsageLimit         *int
	UsedCount          int
	StartsAt           time.Time
	ExpiresAt          time.Time
	IsActive           bool
	CreatedAt          time.Time
	UpdatedAt          time.Time
}
