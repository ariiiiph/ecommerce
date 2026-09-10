package dto

import "time"

type CreateReviewRequest struct {
	ProductID int64   `json:"product_id"`
	Rating    int     `json:"rating"`
	Title     *string `json:"title"`
	Comment   *string `json:"comment"`
}

type UpdateReviewRequest struct {
	Rating  *int    `json:"rating"`
	Title   *string `json:"title"`
	Comment *string `json:"comment"`
}

type ReviewResponse struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	ProductID int64     `json:"product_id"`
	Rating    int       `json:"rating"`
	Title     *string   `json:"title"`
	Comment   *string   `json:"comment"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
