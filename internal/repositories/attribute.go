package repositories

import (
	"context"
	"database/sql"

	"github.com/ariiiiph/ecommerce/internal/models"
)

type AttributeRepository struct {
	db *sql.DB
}

func NewAttributeRepository(db *sql.DB) *AttributeRepository {
	return &AttributeRepository{
		db: db,
	}
}

func (r *AttributeRepository) Create(ctx context.Context, attribute *models.Attribute) error {
	query := `INSERT INTO attributes (name) VALUES ($1) RETURNING id, is_active, created_at, updated_at`

	return r.db.QueryRowContext(ctx, query, attribute.Name).Scan(
		&attribute.ID,
		&attribute.IsActive,
		&attribute.CreatedAt,
		&attribute.UpdatedAt,
	)

}

func (r *AttributeRepository) GetByID(ctx context.Context, id int64) (*models.Attribute, error) {
	query := `SELECT id, name, is_active, created_at, updated_at FROM attributes WHERE id = $1`

	attribute := &models.Attribute{}

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&attribute.ID,
		&attribute.Name,
		&attribute.IsActive,
		&attribute.CreatedAt,
		&attribute.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return attribute, nil
}

func (r *AttributeRepository) GetAll(ctx context.Context) ([]*models.Attribute, error) {
	query := `SELECT id, name, is_active, created_at, updated_at FROM attributes ORDER BY id ASC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var attributes []*models.Attribute

	for rows.Next() {
		attribute := &models.Attribute{}

		err := rows.Scan(
			&attribute.ID,
			&attribute.Name,
			&attribute.IsActive,
			&attribute.CreatedAt,
			&attribute.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		attributes = append(attributes, attribute)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return attributes, nil
}

func (r *AttributeRepository) Update(ctx context.Context, attribute *models.Attribute) error {
	query := `UPDATE attributes SET name = $1, is_active = $2, updated_at = NOW() WHERE id = $3`

	result, err := r.db.ExecContext(
		ctx,
		query,
		attribute.Name,
		attribute.IsActive,
		attribute.ID,
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

func (r *AttributeRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM attributes WHERE id = $1`

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
