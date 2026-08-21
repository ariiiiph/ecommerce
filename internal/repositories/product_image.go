package repositories

import (
	"context"
	"database/sql"

	"github.com/ariiiiph/ecommerce/internal/models"
)

type ProductImageRepository struct {
	db *sql.DB
}

func NewProductImageRepository(db *sql.DB) *ProductImageRepository {
	return &ProductImageRepository{
		db: db,
	}
}

func (r *ProductImageRepository) Create(ctx context.Context, productImage *models.ProductImage) error {
	query := `INSERT INTO product_images (product_id, variant_id, image_url, is_primary, sort_order) VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at`

	return r.db.QueryRowContext(
		ctx,
		query,
		productImage.ProductID,
		productImage.VariantID,
		productImage.ImageURL,
		productImage.IsPrimary,
		productImage.SortOrder,
	).Scan(
		&productImage.ID,
		&productImage.CreatedAt,
	)
}

func (r *ProductImageRepository) GetByID(ctx context.Context, id int64) (*models.ProductImage, error) {
	query := `SELECT id, product_id, variant_id, image_url, is_primary, sort_order, created_at FROM product_images WHERE id = $1`

	productImage := &models.ProductImage{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&productImage.ID,
		&productImage.ProductID,
		&productImage.VariantID,
		&productImage.ImageURL,
		&productImage.IsPrimary,
		&productImage.SortOrder,
		&productImage.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return productImage, nil
}

func (r *ProductImageRepository) GetAllByProductID(ctx context.Context, productID int64) ([]*models.ProductImage, error) {
	query := `SELECT id, product_id, variant_id, image_url, is_primary, sort_order, created_at FROM product_images WHERE product_id = $1 ORDER BY sort_order, id`

	rows, err := r.db.QueryContext(ctx, query, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var productImages []*models.ProductImage
	for rows.Next() {
		productImage := &models.ProductImage{}
		err := rows.Scan(
			&productImage.ID,
			&productImage.ProductID,
			&productImage.VariantID,
			&productImage.ImageURL,
			&productImage.IsPrimary,
			&productImage.SortOrder,
			&productImage.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		productImages = append(productImages, productImage)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return productImages, nil
}

func (r *ProductImageRepository) Update(ctx context.Context, productImage *models.ProductImage) error {
	query := `UPDATE product_images SET variant_id = $1, image_url = $2, is_primary = $3, sort_order = $4 WHERE id = $5`
	result, err := r.db.ExecContext(
		ctx,
		query,
		productImage.VariantID,
		productImage.ImageURL,
		productImage.IsPrimary,
		productImage.SortOrder,
		productImage.ID,
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

func (r *ProductImageRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM product_images WHERE id = $1`

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
