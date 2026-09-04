package dto

import (
	"time"

	"github.com/shopspring/decimal"
)

type CreatePaymentRequest struct {
	OrderID int64 `json:"order_id"`
}

type PaymentResponse struct {
	ID            int64           `json:"id"`
	OrderID       int64           `json:"order_id"`
	Amount        decimal.Decimal `json:"amount"`
	Currency      string          `json:"currency"`
	Provider      string          `json:"provider"`
	TransactionID *string         `json:"transaction_id"`
	Status        string          `json:"status"`
	PaidAt        *time.Time      `json:"paid_at"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
}
