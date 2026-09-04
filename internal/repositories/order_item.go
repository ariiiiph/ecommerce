package repositories

import (
	"context"

	"github.com/ariiiiph/ecommerce/internal/db"
	"github.com/ariiiiph/ecommerce/internal/models"
)

type OrderItemRepository struct {
	db db.DBTX
}

func NewOrderItemRepository(db db.DBTX) *OrderItemRepository {
	return &OrderItemRepository{
		db: db,
	}
}

func (r *OrderItemRepository) Create(ctx context.Context, item *models.OrderItem) error {
	query := `
		INSERT INTO order_items (
			order_id,
			variant_id,
			product_name,
			sku,
			quantity,
			unit_price,
			total_price
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at
	`

	return r.db.QueryRowContext(
		ctx,
		query,
		item.OrderID,
		item.VariantID,
		item.ProductName,
		item.SKU,
		item.Quantity,
		item.UnitPrice,
		item.TotalPrice,
	).Scan(
		&item.ID,
		&item.CreatedAt,
	)
}

func (r *OrderItemRepository) GetByID(ctx context.Context, id int64) (*models.OrderItem, error) {
	query := `
		SELECT
			id,
			order_id,
			variant_id,
			product_name,
			sku,
			quantity,
			unit_price,
			total_price,
			created_at
		FROM order_items
		WHERE id = $1
	`

	item := &models.OrderItem{}

	err := r.db.QueryRowContext(
		ctx,
		query,
		id,
	).Scan(
		&item.ID,
		&item.OrderID,
		&item.VariantID,
		&item.ProductName,
		&item.SKU,
		&item.Quantity,
		&item.UnitPrice,
		&item.TotalPrice,
		&item.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return item, nil
}

func (r *OrderItemRepository) GetAllByOrderID(ctx context.Context, orderID int64) ([]*models.OrderItem, error) {
	query := `
		SELECT
			id,
			order_id,
			variant_id,
			product_name,
			sku,
			quantity,
			unit_price,
			total_price,
			created_at
		FROM order_items
		WHERE order_id = $1
		ORDER BY id ASC
	`

	rows, err := r.db.QueryContext(
		ctx,
		query,
		orderID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]*models.OrderItem, 0)

	for rows.Next() {
		item := &models.OrderItem{}

		err := rows.Scan(
			&item.ID,
			&item.OrderID,
			&item.VariantID,
			&item.ProductName,
			&item.SKU,
			&item.Quantity,
			&item.UnitPrice,
			&item.TotalPrice,
			&item.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (r *OrderItemRepository) DeleteByOrderID(ctx context.Context, orderID int64) error {
	query := `
		DELETE FROM order_items
		WHERE order_id = $1
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		orderID,
	)

	return err
}
