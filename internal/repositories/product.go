package repositories

import (
	"context"
	"database/sql"

	"github.com/ariiiiph/ecommerce/internal/models"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{
		db: db,
	}
}

func (r *ProductRepository) Create(ctx context.Context, product *models.Product) error {
	query := `
		INSERT INTO products (
			name,
			slug,
			description,
			brand_id,
			category_id,
			status
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`

	return r.db.QueryRowContext(
		ctx,
		query,
		product.Name,
		product.Slug,
		product.Description,
		product.BrandID,
		product.CategoryID,
		product.Status,
	).Scan(
		&product.ID,
		&product.CreatedAt,
		&product.UpdatedAt,
	)
}

func (r *ProductRepository) GetByID(ctx context.Context, id int64) (*models.Product, error) {
	query := `
		SELECT
			id,
			name,
			slug,
			description,
			brand_id,
			category_id,
			status,
			created_at,
			updated_at,
			deleted_at
		FROM products
		WHERE id = $1 AND deleted_at IS NULL
	`

	product := &models.Product{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&product.ID,
		&product.Name,
		&product.Slug,
		&product.Description,
		&product.BrandID,
		&product.CategoryID,
		&product.Status,
		&product.CreatedAt,
		&product.UpdatedAt,
		&product.DeletedAt,
	)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (r *ProductRepository) GetAll(ctx context.Context) ([]*models.Product, error) {
	query := `
		SELECT
			id,
			name,
			slug,
			description,
			brand_id,
			category_id,
			status,
			created_at,
			updated_at,
			deleted_at
		FROM products
		WHERE deleted_at IS NULL
		ORDER BY id ASC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*models.Product
	for rows.Next() {
		product := &models.Product{}
		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Slug,
			&product.Description,
			&product.BrandID,
			&product.CategoryID,
			&product.Status,
			&product.CreatedAt,
			&product.UpdatedAt,
			&product.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func (r *ProductRepository) Update(ctx context.Context, product *models.Product) error {
	query := `
		UPDATE products
		SET
			name = $1,
			slug = $2,
			description = $3,
			brand_id = $4,
			category_id = $5,
			status = $6,
			updated_at = NOW()
		WHERE id = $7 AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		product.Name,
		product.Slug,
		product.Description,
		product.BrandID,
		product.CategoryID,
		product.Status,
		product.ID,
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

func (r *ProductRepository) Delete(ctx context.Context, id int64) error {
	query := `
		UPDATE products
		SET deleted_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
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
