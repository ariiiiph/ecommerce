package dto

type CreateInventoryRequest struct {
	VariantID         int64 `json:"variant_id"`
	Quantity          int   `json:"quantity"`
	ReservedQuantity  int   `json:"reserved_quantity"`
	LowStockThreshold int   `json:"low_stock_threshold"`
}

type UpdateInventoryRequest struct {
	Quantity          *int `json:"quantity,omitempty"`
	ReservedQuantity  *int `json:"reserved_quantity,omitempty"`
	LowStockThreshold *int `json:"low_stock_threshold,omitempty"`
}

type InventoryResponse struct {
	VariantID         int64  `json:"variant_id"`
	Quantity          int    `json:"quantity"`
	ReservedQuantity  int    `json:"reserved_quantity"`
	LowStockThreshold int    `json:"low_stock_threshold"`
	CreatedAt         string `json:"created_at"`
	UpdatedAt         string `json:"updated_at"`
}
