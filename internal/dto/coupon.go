package dto

import (
	"time"

	"github.com/shopspring/decimal"
)

type CreateCouponRequest struct {
	Code               string           `json:"code"`
	DiscountType       string           `json:"discount_type"`
	DiscountValue      decimal.Decimal  `json:"discount_value"`
	MinimumOrderAmount *decimal.Decimal `json:"minimum_order_amount"`
	MaximumDiscount    *decimal.Decimal `json:"maximum_discount"`
	UsageLimit         *int             `json:"usage_limit"`
	StartsAt           time.Time        `json:"starts_at"`
	ExpiresAt          time.Time        `json:"expires_at"`
	IsActive           bool             `json:"is_active"`
}

type UpdateCouponRequest struct {
	Code               string           `json:"code"`
	DiscountType       string           `json:"discount_type"`
	DiscountValue      decimal.Decimal  `json:"discount_value"`
	MinimumOrderAmount *decimal.Decimal `json:"minimum_order_amount"`
	MaximumDiscount    *decimal.Decimal `json:"maximum_discount"`
	UsageLimit         *int             `json:"usage_limit"`
	StartsAt           time.Time        `json:"starts_at"`
	ExpiresAt          time.Time        `json:"expires_at"`
	IsActive           bool             `json:"is_active"`
}

type CouponResponse struct {
	ID                 int64            `json:"id"`
	Code               string           `json:"code"`
	DiscountType       string           `json:"discount_type"`
	DiscountValue      decimal.Decimal  `json:"discount_value"`
	MinimumOrderAmount *decimal.Decimal `json:"minimum_order_amount"`
	MaximumDiscount    *decimal.Decimal `json:"maximum_discount"`
	UsageLimit         *int             `json:"usage_limit"`
	UsedCount          int              `json:"used_count"`
	StartsAt           time.Time        `json:"starts_at"`
	ExpiresAt          time.Time        `json:"expires_at"`
	IsActive           bool             `json:"is_active"`
	CreatedAt          time.Time        `json:"created_at"`
	UpdatedAt          time.Time        `json:"updated_at"`
}
