package repositories

import (
	"context"
	"database/sql"

	"github.com/ariiiiph/ecommerce/internal/models"
)

type WishlistItemRepository struct {
	db *sql.DB
}

func NewWishlistItemRepository(db *sql.DB) *WishlistItemRepository {
	return &WishlistItemRepository{
		db: db,
	}
}

func (r *WishlistItemRepository) Create(ctx context.Context, item *models.WishlistItem) error {
	query := `
		INSERT INTO wishlist_items (wishlist_id, product_id)
		VALUES ($1, $2)
		RETURNING id, wishlist_id, product_id, created_at
	`

	return r.db.QueryRowContext(
		ctx,
		query,
		item.WishlistID,
		item.ProductID,
	).Scan(
		&item.ID,
		&item.WishlistID,
		&item.ProductID,
		&item.CreatedAt,
	)
}

func (r *WishlistItemRepository) GetByID(ctx context.Context, id int64) (*models.WishlistItem, error) {
	query := `
		SELECT id, wishlist_id, product_id, created_at
		FROM wishlist_items
		WHERE id = $1
	`

	item := &models.WishlistItem{}

	err := r.db.QueryRowContext(
		ctx,
		query,
		id,
	).Scan(
		&item.ID,
		&item.WishlistID,
		&item.ProductID,
		&item.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return item, nil
}

func (r *WishlistItemRepository) GetAllByWishlistID(ctx context.Context, wishlistID int64) ([]*models.WishlistItem, error) {
	query := `
		SELECT id, wishlist_id, product_id, created_at
		FROM wishlist_items
		WHERE wishlist_id = $1
		ORDER BY id ASC
	`

	rows, err := r.db.QueryContext(ctx, query, wishlistID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var items []*models.WishlistItem

	for rows.Next() {
		item := &models.WishlistItem{}

		if err := rows.Scan(
			&item.ID,
			&item.WishlistID,
			&item.ProductID,
			&item.CreatedAt,
		); err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (r *WishlistItemRepository) Delete(ctx context.Context, id int64) error {
	query := `
		DELETE FROM wishlist_items
		WHERE id = $1
	`

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
