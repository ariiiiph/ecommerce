package models

import "time"

type Product struct {
	ID          int64
	Name        string
	Slug        string
	Description string
	BrandID     int64
	CategoryID  int64
	Status      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}
