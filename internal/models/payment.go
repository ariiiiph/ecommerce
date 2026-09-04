package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type Payment struct {
	ID            int64
	OrderID       int64
	Amount        decimal.Decimal
	Currency      string
	Provider      string
	TransactionID *string
	Status        string
	PaidAt        *time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
