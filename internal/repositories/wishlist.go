package repositories

import (
	"context"
	"database/sql"

	"github.com/ariiiiph/ecommerce/internal/models"
)

type WishlistRepository struct {
	db *sql.DB
}

func NewWishlistRepository(db *sql.DB) *WishlistRepository {
	return &WishlistRepository{
		db: db,
	}
}

func (r *WishlistRepository) Create(ctx context.Context, wishlist *models.Wishlist) error {
	query := `INSERT INTO wishlists (user_id) VALUES ($1) RETURNING id, user_id, created_at`

	return r.db.QueryRowContext(
		ctx,
		query,
		wishlist.UserID,
	).Scan(
		&wishlist.ID,
		&wishlist.UserID,
		&wishlist.CreatedAt,
	)
}

func (r *WishlistRepository) GetByID(ctx context.Context, id int64) (*models.Wishlist, error) {
	query := `SELECT id, user_id, created_at FROM wishlists WHERE id = $1`

	wishlist := &models.Wishlist{}

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&wishlist.ID,
		&wishlist.UserID,
		&wishlist.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return wishlist, nil
}

func (r *WishlistRepository) GetByUserID(ctx context.Context, userID int64) (*models.Wishlist, error) {
	query := `SELECT id, user_id, created_at FROM wishlists WHERE user_id = $1`

	wishlist := &models.Wishlist{}

	err := r.db.QueryRowContext(
		ctx,
		query,
		userID,
	).Scan(
		&wishlist.ID,
		&wishlist.UserID,
		&wishlist.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return wishlist, nil
}

func (r *WishlistRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM wishlists WHERE id = $1`

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
