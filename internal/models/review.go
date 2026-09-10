package models

import "time"

type Review struct {
	ID        int64
	UserID    int64
	ProductID int64
	Rating    int
	Title     *string
	Comment   *string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}
