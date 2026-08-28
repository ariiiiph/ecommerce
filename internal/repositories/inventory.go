package repositories

import (
	"context"
	"database/sql"

	"github.com/ariiiiph/ecommerce/internal/models"
)

type InventoryRepository struct {
	db *sql.DB
}

func NewInventoryRepository(db *sql.DB) *InventoryRepository {
	return &InventoryRepository{
		db: db,
	}
}

func (r *InventoryRepository) Create(ctx context.Context, inventory *models.Inventory) error {
	query := `INSERT INTO inventory (variant_id, quantity, reserved_quantity, low_stock_threshold) VALUES ($1, $2, $3, $4) RETURNING variant_id, quantity, reserved_quantity, low_stock_threshold, created_at, updated_at`

	return r.db.QueryRowContext(
		ctx,
		query,
		inventory.VariantID,
		inventory.Quantity,
		inventory.ReservedQuantity,
		inventory.LowStockThreshold,
	).Scan(
		&inventory.VariantID,
		&inventory.Quantity,
		&inventory.ReservedQuantity,
		&inventory.LowStockThreshold,
		&inventory.CreatedAt,
		&inventory.UpdatedAt,
	)
}

func (r *InventoryRepository) GetByVariantID(ctx context.Context, variantID int64) (*models.Inventory, error) {
	query := `SELECT variant_id, quantity, reserved_quantity, low_stock_threshold, created_at, updated_at FROM inventory WHERE variant_id = $1`

	inventory := &models.Inventory{}

	err := r.db.QueryRowContext(
		ctx,
		query,
		variantID,
	).Scan(
		&inventory.VariantID,
		&inventory.Quantity,
		&inventory.ReservedQuantity,
		&inventory.LowStockThreshold,
		&inventory.CreatedAt,
		&inventory.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return inventory, nil
}

func (r *InventoryRepository) GetAll(ctx context.Context) ([]*models.Inventory, error) {
	query := `SELECT variant_id, quantity, reserved_quantity, low_stock_threshold, created_at, updated_at FROM inventory ORDER BY variant_id ASC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var inventories []*models.Inventory

	for rows.Next() {
		inventory := &models.Inventory{}

		if err := rows.Scan(
			&inventory.VariantID,
			&inventory.Quantity,
			&inventory.ReservedQuantity,
			&inventory.LowStockThreshold,
			&inventory.CreatedAt,
			&inventory.UpdatedAt,
		); err != nil {
			return nil, err
		}
		inventories = append(inventories, inventory)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return inventories, nil
}

func (r *InventoryRepository) Update(ctx context.Context, inventory *models.Inventory) error {
	query := `UPDATE inventory SET quantity = $1, reserved_quantity = $2, low_stock_threshold = $3, updated_at = NOW() WHERE variant_id = $4`

	result, err := r.db.ExecContext(
		ctx,
		query,
		inventory.Quantity,
		inventory.ReservedQuantity,
		inventory.LowStockThreshold,
		inventory.VariantID,
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

func (r *InventoryRepository) Delete(ctx context.Context, variantID int64) error {
	query := `DELETE FROM inventory WHERE variant_id = $1`

	result, err := r.db.ExecContext(ctx, query, variantID)
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
