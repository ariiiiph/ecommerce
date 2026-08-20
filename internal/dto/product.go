package dto

type CreateProductRequest struct {
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	BrandID     int64  `json:"brand_id"`
	CategoryID  int64  `json:"category_id"`
	Status      string `json:"status"`
}

type UpdateProductRequest struct {
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	BrandID     int64  `json:"brand_id"`
	CategoryID  int64  `json:"category_id"`
	Status      string `json:"status"`
}

type ProductResponse struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	BrandID     int64  `json:"brand_id"`
	CategoryID  int64  `json:"category_id"`
	Status      string `json:"status"`
}
