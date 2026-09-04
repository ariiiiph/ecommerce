package repositories

import (
	"context"
	"database/sql"

	"github.com/ariiiiph/ecommerce/internal/db"
	"github.com/ariiiiph/ecommerce/internal/models"
)

type OrderRepository struct {
	db db.DBTX
}

func NewOrderRepository(db db.DBTX) *OrderRepository {
	return &OrderRepository{
		db: db,
	}
}

func (r *OrderRepository) Create(ctx context.Context, order *models.Order) error {
	query := `
		INSERT INTO orders (
			order_number,
			user_id,
			address_id,
			coupon_id,
			status,
			subtotal,
			discount_amount,
			shipping_amount,
			total_amount,
			shipping_recipient_name,
			shipping_phone,
			shipping_country,
			shipping_city,
			shipping_address_line,
			shipping_postal_code
		)
		VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8, $9,$10,
			$11, $12, $13, $14, $15
		)
		RETURNING
			id,
			created_at,
			updated_at
	`
	return r.db.QueryRowContext(
		ctx,
		query,
		order.OrderNumber,
		order.UserID,
		order.AddressID,
		order.CouponID,
		order.Status,
		order.Subtotal,
		order.DiscountAmount,
		order.ShippingAmount,
		order.TotalAmount,
		order.ShippingRecipientName,
		order.ShippingPhone,
		order.ShippingCountry,
		order.ShippingCity,
		order.ShippingAddressLine,
		order.ShippingPostalCode,
	).Scan(
		&order.ID,
		&order.CreatedAt,
		&order.UpdatedAt,
	)
}

func (r *OrderRepository) GetByID(ctx context.Context, id int64) (*models.Order, error) {
	query := `
		SELECT
			id,
			order_number,
			user_id,
			address_id,
			coupon_id,
			status,
			subtotal,
			discount_amount,
			shipping_amount,
			total_amount,
			shipping_recipient_name,
			shipping_phone,
			shipping_country,
			shipping_city,
			shipping_address_line,
			shipping_postal_code,
			created_at,
			updated_at
		FROM orders
		WHERE id = $1
	`
	order := &models.Order{}

	err := r.db.QueryRowContext(
		ctx,
		query,
		id,
	).Scan(
		&order.ID,
		&order.OrderNumber,
		&order.UserID,
		&order.AddressID,
		&order.CouponID,
		&order.Status,
		&order.Subtotal,
		&order.DiscountAmount,
		&order.ShippingAmount,
		&order.TotalAmount,
		&order.ShippingRecipientName,
		&order.ShippingPhone,
		&order.ShippingCountry,
		&order.ShippingCity,
		&order.ShippingAddressLine,
		&order.ShippingPostalCode,
		&order.CreatedAt,
		&order.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return order, nil

}

func (r *OrderRepository) GetAllByUserID(ctx context.Context, userID int64) ([]*models.Order, error) {
	query := `
		SELECT
			id,
			order_number,
			user_id,
			address_id,
			coupon_id,
			status,
			subtotal,
			discount_amount,
			shipping_amount,
			total_amount,
			shipping_recipient_name,
			shipping_phone,
			shipping_country,
			shipping_city,
			shipping_address_line,
			shipping_postal_code,
			created_at,
			updated_at
		FROM orders
		WHERE user_id = $1
		ORDER BY id DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	orders := make([]*models.Order, 0)

	for rows.Next() {
		order := &models.Order{}

		err := rows.Scan(
			&order.ID,
			&order.OrderNumber,
			&order.UserID,
			&order.AddressID,
			&order.CouponID,
			&order.Status,
			&order.Subtotal,
			&order.DiscountAmount,
			&order.ShippingAmount,
			&order.TotalAmount,
			&order.ShippingRecipientName,
			&order.ShippingPhone,
			&order.ShippingCountry,
			&order.ShippingCity,
			&order.ShippingAddressLine,
			&order.ShippingPostalCode,
			&order.CreatedAt,
			&order.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return orders, nil

}

func (r *OrderRepository) UpdateStatus(ctx context.Context, id int64, status string) error {
	query := `UPDATE orders SET status = $1, updated_at = NOW() WHERE id = $2`

	result, err := r.db.ExecContext(ctx, query, status, id)

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

func (r *OrderRepository) Delete(ctx context.Context, id int64) error {
	query := `
		DELETE FROM orders
		WHERE id = $1
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		id,
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
