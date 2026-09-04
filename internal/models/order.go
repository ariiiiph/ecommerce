package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type Order struct {
	ID int64

	OrderNumber string

	UserID    int64
	AddressID *int64
	CouponID  *int64

	Status string

	Subtotal       decimal.Decimal
	DiscountAmount decimal.Decimal
	ShippingAmount decimal.Decimal
	TotalAmount    decimal.Decimal

	ShippingRecipientName string
	ShippingPhone         string
	ShippingCountry       string
	ShippingAddressLine   string
	ShippingCity          string

	ShippingPostalCode string

	CreatedAt time.Time
	UpdatedAt time.Time
}
