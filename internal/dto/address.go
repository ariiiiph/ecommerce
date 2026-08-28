package dto

import "time"

type CreateAddressRequest struct {
	UserID        int64  `json:"user_id"`
	Title         string `json:"title"`
	RecipientName string `json:"recipient_name"`
	Phone         string `json:"phone"`
	Country       string `json:"country"`
	City          string `json:"city"`
	AddressLine   string `json:"address_line"`
	PostalCode    string `json:"postal_code"`
	IsDefault     bool   `json:"is_default"`
}

type UpdateAddressRequest struct {
	Title         *string `json:"title"`
	RecipientName *string `json:"recipient_name"`
	Phone         *string `json:"phone"`
	Country       *string `json:"country"`
	City          *string `json:"city"`
	AddressLine   *string `json:"address_line"`
	PostalCode    *string `json:"postal_code"`
	IsDefault     *bool   `json:"is_default"`
}

type AddressResponse struct {
	ID            int64     `json:"id"`
	UserID        int64     `json:"user_id"`
	Title         string    `json:"title"`
	RecipientName string    `json:"recipient_name"`
	Phone         string    `json:"phone"`
	Country       string    `json:"country"`
	City          string    `json:"city"`
	AddressLine   string    `json:"address_line"`
	PostalCode    string    `json:"postal_code"`
	IsDefault     bool      `json:"is_default"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
