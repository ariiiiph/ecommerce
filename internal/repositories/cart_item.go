package repositories

import (
	"context"
	"database/sql"

	"github.com/ariiiiph/ecommerce/internal/models"
)

type CartItemRepository struct {
	db *sql.DB
}

func NewCartItemRepository(db *sql.DB) *CartItemRepository {
	return &CartItemRepository{
		db: db,
	}
}

func (r *CartItemRepository) Create(ctx context.Context, cartItem *models.CartItem) error {
	query := `INSERT INTO cart_items (cart_id, variant_id, quantity) VALUES ($1, $2, $3) RETURNING id, cart_id, variant_id, quantity, created_at, updated_at`

	return r.db.QueryRowContext(
		ctx,
		query,
		cartItem.CartID,
		cartItem.VariantID,
		cartItem.Quantity,
	).Scan(
		&cartItem.ID,
		&cartItem.CartID,
		&cartItem.VariantID,
		&cartItem.Quantity,
		&cartItem.CreatedAt,
		&cartItem.UpdatedAt,
	)

}

func (r *CartItemRepository) GetByID(ctx context.Context, id int64) (*models.CartItem, error) {
	query := `SELECT id, cart_id, variant_id, Quantity, created_at, updated_at FROM cart_items WHERE id = $1`

	cartItem := &models.CartItem{}

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&cartItem.ID,
		&cartItem.CartID,
		&cartItem.VariantID,
		&cartItem.Quantity,
		&cartItem.CreatedAt,
		&cartItem.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return cartItem, nil
}

func (r *CartItemRepository) GetAllByCartID(ctx context.Context, cartID int64) ([]*models.CartItem, error) {
	query := `SELECT id, cart_id, variant_id, Quantity, created_at, updated_at FROM cart_items WHERE cart_id = $1 ORDER BY id ASC`

	rows, err := r.db.QueryContext(ctx, query, cartID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var cartItems []*models.CartItem

	for rows.Next() {
		cartItem := &models.CartItem{}

		err := rows.Scan(
			&cartItem.ID,
			&cartItem.CartID,
			&cartItem.VariantID,
			&cartItem.Quantity,
			&cartItem.CreatedAt,
			&cartItem.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		cartItems = append(cartItems, cartItem)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return cartItems, nil
}

func (r *CartItemRepository) Update(ctx context.Context, cartItem *models.CartItem) error {
	query := `UPDATE cart_items SET quantity = $1, updated_at = NOW() WHERE id = $2`

	result, err := r.db.ExecContext(
		ctx,
		query,
		cartItem.Quantity,
		cartItem.ID,
	)
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

func (r *CartItemRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM cart_items WHERE id = $1`

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
