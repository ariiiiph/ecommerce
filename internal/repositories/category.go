package repositories

import (
	"context"
	"database/sql"

	"github.com/ariiiiph/ecommerce/internal/models"
)

type CategoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{
		db: db,
	}
}

func (r *CategoryRepository) Create(ctx context.Context, category *models.Category) error {
	query := `INSERT INTO categories(name, slug, description, parent_id) VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at`

	return r.db.QueryRowContext(
		ctx,
		query,
		category.Name,
		category.Slug,
		category.Description,
		category.ParentID,
	).Scan(
		&category.ID,
		&category.CreatedAt,
		&category.UpdatedAt,
	)
}

func (r *CategoryRepository) FindByID(ctx context.Context, id int64) (*models.Category, error) {
	query := `SELECT id, name, slug, description,parent_id, created_at, updated_at FROM categories WHERE id = $1`

	category := &models.Category{}

	err := r.db.QueryRowContext(
		ctx,
		query,
		id,
	).Scan(
		&category.ID,
		&category.Name,
		&category.Slug,
		&category.Description,
		&category.ParentID,
		&category.CreatedAt,
		&category.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return category, nil
}

func (r *CategoryRepository) FindAll(ctx context.Context) ([]*models.Category, error) {
	query := `SELECT id, name, slug, description,parent_id, created_at, updated_at FROM categories ORDER BY id ASC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*models.Category
	for rows.Next() {
		category := &models.Category{}
		err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.Slug,
			&category.Description,
			&category.ParentID,
			&category.CreatedAt,
			&category.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *CategoryRepository) Update(ctx context.Context, category *models.Category) error {
	query := `UPDATE categories SET name = $1, slug = $2, description = $3, parent_id = $4, updated_at = $5 WHERE id = $6`

	_, err := r.db.ExecContext(
		ctx,
		query,
		category.Name,
		category.Slug,
		category.Description,
		category.ParentID,
		category.UpdatedAt,
		category.ID,
	)
	return err
}

func (r *CategoryRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM categories WHERE id = $1`

	_, err := r.db.ExecContext(
		ctx,
		query,
		id,
	)
	return err
}

func (r *CategoryRepository) FindChildren(ctx context.Context, parentID int64) ([]*models.Category, error) {
	query := `SELECT id, name, slug, description,parent_id, created_at, updated_at FROM categories WHERE parent_id = $1`

	rows, err := r.db.QueryContext(ctx, query, parentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*models.Category
	for rows.Next() {
		category := &models.Category{}
		err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.Slug,
			&category.Description,
			&category.ParentID,
			&category.CreatedAt,
			&category.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *CategoryRepository) FindTree(ctx context.Context, parentID int64) ([]*models.Category, error) {
	query := `WITH RECURSIVE category_tree AS (
		SELECT id, name, slug, description, parent_id, created_at, updated_at
		FROM categories
		WHERE parent_id = $1
		UNION ALL
		SELECT c.id, c.name, c.slug, c.description, c.parent_id, c.created_at, c.updated_at
		FROM categories c
		INNER JOIN category_tree ct ON ct.id = c.parent_id
	)
	SELECT id, name, slug, description, parent_id, created_at, updated_at FROM category_tree ORDER BY id ASC`
	rows, err := r.db.QueryContext(ctx, query, parentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*models.Category
	for rows.Next() {
		category := &models.Category{}
		err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.Slug,
			&category.Description,
			&category.ParentID,
			&category.CreatedAt,
			&category.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}
