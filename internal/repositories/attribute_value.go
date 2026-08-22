package repositories

import (
	"context"
	"database/sql"

	"github.com/ariiiiph/ecommerce/internal/models"
)

type AttributeValueRepository struct {
	db *sql.DB
}

func NewAttributeValueRepository(db *sql.DB) *AttributeValueRepository {
	return &AttributeValueRepository{
		db: db,
	}
}

func (r *AttributeValueRepository) Create(ctx context.Context, attributeValue *models.AttributeValue) error {
	query := `INSERT INTO attribute_values (attribute_id, value) VALUES ($1, $2) RETURNING id, attribute_id, value, created_at`

	return r.db.QueryRowContext(ctx, query, attributeValue.AttributeID, attributeValue.Value).Scan(
		&attributeValue.ID,
		&attributeValue.AttributeID,
		&attributeValue.Value,
		&attributeValue.CreatedAt,
	)
}

func (r *AttributeValueRepository) GetByID(ctx context.Context, id int64) (*models.AttributeValue, error) {
	query := `SELECT id, attribute_id, value, created_at FROM attribute_values WHERE id = $1`

	attributeValue := &models.AttributeValue{}

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&attributeValue.ID,
		&attributeValue.AttributeID,
		&attributeValue.Value,
		&attributeValue.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return attributeValue, nil
}

func (r *AttributeValueRepository) GetAll(ctx context.Context) ([]*models.AttributeValue, error) {
	query := `SELECT id, attribute_id, value, created_at FROM attribute_values ORDER BY id ASC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var attributeValues []*models.AttributeValue

	for rows.Next() {
		attributeValue := &models.AttributeValue{}

		err := rows.Scan(
			&attributeValue.ID,
			&attributeValue.AttributeID,
			&attributeValue.Value,
			&attributeValue.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		attributeValues = append(attributeValues, attributeValue)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return attributeValues, nil

}

func (r *AttributeValueRepository) Update(ctx context.Context, attributeValue *models.AttributeValue) error {
	query := `UPDATE attribute_values SET value = $1 WHERE id = $2`

	result, err := r.db.ExecContext(ctx, query, attributeValue.Value, attributeValue.ID)

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

func (r *AttributeValueRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM attribute_values WHERE id = $1`

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
