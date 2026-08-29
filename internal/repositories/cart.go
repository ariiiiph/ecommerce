package repositories

import (
	"context"
	"database/sql"

	"github.com/ariiiiph/ecommerce/internal/models"
)

type CartRepository struct {
	db *sql.DB
}

func NewCartRepository(db *sql.DB) *CartRepository {
	return &CartRepository{
		db: db,
	}
}

func (r *CartRepository) Create(ctx context.Context, cart *models.Cart) error {
	query := `INSERT INTO carts (user_id) VALUES ($1) RETURNING id, user_id, created_at, updated_at`

	return r.db.QueryRowContext(ctx, query, cart.UserID).Scan(
		&cart.ID,
		&cart.UserID,
		&cart.CreatedAt,
		&cart.UpdatedAt,
	)
}

func (r *CartRepository) GetByID(ctx context.Context, id int64) (*models.Cart, error) {
	query := `SELECT id, user_id, created_at, updated_at FROM carts WHERE id = $1`

	cart := &models.Cart{}

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&cart.ID,
		&cart.UserID,
		&cart.CreatedAt,
		&cart.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return cart, nil
}

func (r *CartRepository) GetByUserID(ctx context.Context, userID int64) (*models.Cart, error) {
	query := `SELECT id, user_id, created_at, updated_at FROM carts WHERE user_id = $1`

	cart := &models.Cart{}

	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&cart.ID,
		&cart.UserID,
		&cart.CreatedAt,
		&cart.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return cart, nil

}

func (r *CartRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM carts WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil

}
