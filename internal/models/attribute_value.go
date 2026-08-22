package models

import "time"

type AttributeValue struct {
	ID          int64     `json:"id"`
	AttributeID int64     `json:"attribute_id"`
	Value       string    `json:"value"`
	CreatedAt   time.Time `json:"created_at"`
}
