package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type OrderItem struct {
	ID int64

	OrderID   int64
	VariantID int64

	ProductName string
	SKU         string

	Quantity int

	UnitPrice  decimal.Decimal
	TotalPrice decimal.Decimal

	CreatedAt time.Time
}
