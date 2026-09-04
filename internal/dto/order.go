package dto

import (
	"time"

	"github.com/shopspring/decimal"
)

type CreateOrderRequest struct {
	AddressID  *int64  `json:"address_id"`
	CouponCode *string `json:"coupon_id"`
}

type UpdateOrderStatusRequest struct {
	Status string `json:"status"`
}

type OrderItemResponse struct {
	ID int64 `json:"id"`

	OrderID   int64 `json:"order_id"`
	VariantID int64 `json:"variant_id"`

	ProductName string `json:"product_name"`
	SKU         string `json:"sku"`

	Quantity int `json:"quantity"`

	UnitPrice  decimal.Decimal `json:"unit_price"`
	TotalPrice decimal.Decimal `json:"total_price"`

	CreatedAt time.Time `json:"created_at"`
}

type OrderResponse struct {
	ID          int64  `json:"id"`
	OrderNumber string `json:"order_number"`

	UserID    int64  `json:"user_id"`
	AddressID *int64 `json:"address_id"`
	CouponID  *int64 `json:"coupon_id"`

	Status string `json:"status"`

	Subtotal       decimal.Decimal `json:"subtotal"`
	DiscountAmount decimal.Decimal `json:"discount_amount"`
	ShippingAmount decimal.Decimal `json:"shipping_amount"`
	TotalAmount    decimal.Decimal `json:"total_amount"`

	ShippingRecipientName string `json:"shipping_recipient_name"`
	ShippingPhone         string `json:"shipping_phone"`
	ShippingCountry       string `json:"shipping_country"`
	ShippingCity          string `json:"shipping_city"`
	ShippingAddressLine   string `json:"shipping_address_line"`
	ShippingPostalCode    string `json:"shipping_postal_code"`

	Items []OrderItemResponse `json:"items"`

	CreatedAt  time.Time `json:"created_at"`
	UpdatedtAt time.Time `json:"updated_at"`
}
