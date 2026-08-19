package repositories

import (
	"context"
	"database/sql"

	"github.com/ariiiiph/ecommerce/internal/models"
)

type BrandRepository struct {
	db *sql.DB
}

func NewBrandRepository(db *sql.DB) *BrandRepository {
	return &BrandRepository{
		db: db,
	}
}

func (r *BrandRepository) Create(ctx context.Context, brand *models.Brand) error {
	query := `INSERT INTO brands(name, slug, description) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at`

	return r.db.QueryRowContext(
		ctx,
		query,
		brand.Name,
		brand.Slug,
		brand.Description,
	).Scan(
		&brand.ID,
		&brand.CreatedAt,
		&brand.UpdatedAt,
	)
}

func (r *BrandRepository) FindByID(ctx context.Context, id int64) (*models.Brand, error) {
	query := `SELECT id, name, slug, description, created_at, updated_at FROM brands WHERE id = $1`

	brand := &models.Brand{}

	err := r.db.QueryRowContext(
		ctx,
		query,
		id,
	).Scan(
		&brand.ID,
		&brand.Name,
		&brand.Slug,
		&brand.Description,
		&brand.CreatedAt,
		&brand.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return brand, nil
}

func (r *BrandRepository) FindAll(ctx context.Context) ([]*models.Brand, error) {
	query := `SELECT id, name, slug, description, created_at, updated_at FROM brands ORDER BY id ASC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var brands []*models.Brand
	for rows.Next() {
		brand := &models.Brand{}
		err := rows.Scan(
			&brand.ID,
			&brand.Name,
			&brand.Slug,
			&brand.Description,
			&brand.CreatedAt,
			&brand.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		brands = append(brands, brand)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return brands, nil

}

func (r *BrandRepository) Update(ctx context.Context, brand *models.Brand) error {
	query := `UPDATE brands SET name = $1, slug = $2, description = $3, updated_at = $4 WHERE id = $5`

	_, err := r.db.ExecContext(
		ctx,
		query,
		brand.Name,
		brand.Slug,
		brand.Description,
		brand.UpdatedAt,
		brand.ID,
	)
	return err
}

func (r *BrandRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM brands WHERE id = $1`

	_, err := r.db.ExecContext(
		ctx,
		query,
		id,
	)
	return err
}
