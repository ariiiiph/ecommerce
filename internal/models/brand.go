package models

import "time"

type Brand struct {
	ID          int64
	Name        string
	Slug        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
