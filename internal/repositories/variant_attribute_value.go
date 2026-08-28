package repositories

import (
	"context"
	"database/sql"

	"github.com/ariiiiph/ecommerce/internal/models"
)

type VariantAttributeValueRepository struct {
	db *sql.DB
}

func NewVariantAttributeValueRepository(db *sql.DB) *VariantAttributeValueRepository {
	return &VariantAttributeValueRepository{
		db: db,
	}
}

func (r *VariantAttributeValueRepository) Create(ctx context.Context, relationship *models.VariantAttributeValue) error {
	query := `
		INSERT INTO variant_attribute_values (variant_id, attribute_value_id) VALUES ($1, $2)`

	_, err := r.db.ExecContext(
		ctx,
		query,
		relationship.VariantID,
		relationship.AttributeValueID,
	)

	return err
}

func (r *VariantAttributeValueRepository) GetByVariantID(ctx context.Context, variantID int64) ([]*models.VariantAttributeValue, error) {
	query := `SELECT variant_id,attribute_value_id FROM variant_attribute_values WHERE variant_id = $1 ORDER BY attribute_value_id ASC`

	rows, err := r.db.QueryContext(ctx, query, variantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var relationships []*models.VariantAttributeValue

	for rows.Next() {
		relationship := &models.VariantAttributeValue{}

		if err := rows.Scan(
			&relationship.VariantID,
			&relationship.AttributeValueID,
		); err != nil {
			return nil, err
		}

		relationships = append(relationships, relationship)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return relationships, nil
}

func (r *VariantAttributeValueRepository) Delete(ctx context.Context, variantID int64, attributeValueID int64) error {
	query := `DELETE FROM variant_attribute_valuesWHERE variant_id = $1AND attribute_value_id = $2`

	result, err := r.db.ExecContext(
		ctx,
		query,
		variantID,
		attributeValueID,
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
