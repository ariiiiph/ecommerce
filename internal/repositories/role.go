package repositories

import (
	"context"
	"database/sql"

	"github.com/ariiiiph/ecommerce/internal/models"
)

type RoleRepository struct {
	db *sql.DB
}

func NewRoleRepository(db *sql.DB) *RoleRepository {
	return &RoleRepository{
		db: db,
	}
}

func (r *RoleRepository) FindByID(ctx context.Context, id int64) (*models.Role, error) {
	query := `
		SELECT id, name
		FROM roles
		WHERE id = $1
	`

	role := &models.Role{}

	err := r.db.QueryRowContext(
		ctx,
		query,
		id,
	).Scan(
		&role.ID,
		&role.Name,
	)
	if err != nil {
		return nil, err
	}

	return role, nil
}

func (r *RoleRepository) FindByName(ctx context.Context, name string) (*models.Role, error) {
	query := `
		SELECT id, name
		FROM roles
		WHERE name = $1
	`

	role := &models.Role{}

	err := r.db.QueryRowContext(
		ctx,
		query,
		name,
	).Scan(
		&role.ID,
		&role.Name,
	)
	if err != nil {
		return nil, err
	}

	return role, nil
}
