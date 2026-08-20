package dto

type CreateProductVariantRequest struct {
	ProductID     int64    `json:"product_id"`
	SKU           string   `json:"sku"`
	Price         float64  `json:"price"`
	DiscountPrice *float64 `json:"discount_price,omitempty"`
}

type UpdateProductVariantRequest struct {
	SKU           *string  `json:"sku,omitempty"`
	Price         *float64 `json:"price,omitempty"`
	DiscountPrice *float64 `json:"discount_price,omitempty"`
}

type ProductVariantResponse struct {
	ID            int64    `json:"id"`
	ProductID     int64    `json:"product_id"`
	SKU           string   `json:"sku"`
	Price         float64  `json:"price"`
	DiscountPrice *float64 `json:"discount_price,omitempty"`
}
