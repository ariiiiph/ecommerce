package dto

type CreateProductImageRequest struct {
	ProductID int64  `json:"product_id"`
	VariantID *int64 `json:"variant_id,omitempty"`
	ImageURL  string `json:"image_url"`
	IsPrimary bool   `json:"is_primary"`
	SortOrder int    `json:"sort_order"`
}

type UpdateProductImageRequest struct {
	VariantID *int64  `json:"variant_id,omitempty"`
	ImageURL  *string `json:"image_url,omitempty"`
	IsPrimary *bool   `json:"is_primary,omitempty"`
	SortOrder *int    `json:"sort_order,omitempty"`
}

type ProductImageResponse struct {
	ID        int64  `json:"id"`
	ProductID int64  `json:"product_id"`
	VariantID *int64 `json:"variant_id,omitempty"`
	ImageURL  string `json:"image_url"`
	IsPrimary bool   `json:"is_primary"`
	SortOrder int    `json:"sort_order"`
}
