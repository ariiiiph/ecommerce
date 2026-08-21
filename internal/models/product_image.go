package models

import "time"

type ProductImage struct {
	ID        int64
	ProductID int64
	VariantID *int64
	ImageURL  string
	IsPrimary bool
	SortOrder int
	CreatedAt time.Time
}
