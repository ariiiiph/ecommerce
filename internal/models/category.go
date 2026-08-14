package models

import "time"

type Category struct {
	ID          int64
	Name        string
	Slug        string
	Description string
	ParentID    *int64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
