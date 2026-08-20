package repositories

import (
	"context"
	"database/sql"

	"github.com/ariiiiph/ecommerce/internal/models"
)

type ProductVariantRepository struct {
	db *sql.DB
}

func NewProductVariantRepository(db *sql.DB) *ProductVariantRepository {
	return &ProductVariantRepository{
		db: db,
	}
}

func (r *ProductVariantRepository) Create(ctx context.Context, productVariant *models.ProductVariant) error {
	query := `INSERT INTO product_variants (product_id, sku, price, discount_price) VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at`

	return r.db.QueryRowContext(
		ctx,
		query,
		productVariant.ProductID,
		productVariant.SKU,
		productVariant.Price,
		productVariant.DiscountPrice,
	).Scan(
		&productVariant.ID,
		&productVariant.CreatedAt,
		&productVariant.UpdatedAt,
	)
}

func (r *ProductVariantRepository) GetByID(ctx context.Context, id int64) (*models.ProductVariant, error) {
	query := `SELECT id, product_id, sku, price, discount_price, created_at, updated_at FROM product_variants WHERE id = $1`

	productVariant := &models.ProductVariant{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&productVariant.ID,
		&productVariant.ProductID,
		&productVariant.SKU,
		&productVariant.Price,
		&productVariant.DiscountPrice,
		&productVariant.CreatedAt,
		&productVariant.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return productVariant, nil
}

func (r *ProductVariantRepository) GetAllByProductID(ctx context.Context, productID int64) ([]*models.ProductVariant, error) {
	query := `SELECT id, product_id, sku, price, discount_price, created_at, updated_at FROM product_variants WHERE product_id = $1 ORDER BY id`

	rows, err := r.db.QueryContext(ctx, query, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var productVariants []*models.ProductVariant
	for rows.Next() {
		productVariant := &models.ProductVariant{}
		err := rows.Scan(
			&productVariant.ID,
			&productVariant.ProductID,
			&productVariant.SKU,
			&productVariant.Price,
			&productVariant.DiscountPrice,
			&productVariant.CreatedAt,
			&productVariant.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		productVariants = append(productVariants, productVariant)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return productVariants, nil
}

func (r *ProductVariantRepository) Update(ctx context.Context, productVariant *models.ProductVariant) error {
	query := `UPDATE product_variants SET sku = $1, price = $2, discount_price = $3, updated_at = CURRENT_TIMESTAMP WHERE id = $4`
	result, err := r.db.ExecContext(
		ctx,
		query,
		productVariant.SKU,
		productVariant.Price,
		productVariant.DiscountPrice,
		productVariant.ID,
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

func (r *ProductVariantRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM product_variants WHERE id = $1`
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
