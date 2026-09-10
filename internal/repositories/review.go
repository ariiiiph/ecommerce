package repositories

import (
	"context"

	"github.com/ariiiiph/ecommerce/internal/db"
	"github.com/ariiiiph/ecommerce/internal/models"
)

type ReviewRepository struct {
	db db.DBTX
}

func NewReviewRepository(db db.DBTX) *ReviewRepository {
	return &ReviewRepository{
		db: db,
	}
}

func (r *ReviewRepository) Create(ctx context.Context, review *models.Review) error {
	query := `
		INSERT INTO reviews (
			user_id,
			product_id,
			rating,
			title,
			comment,
			status
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`

	return r.db.QueryRowContext(
		ctx,
		query,
		review.UserID,
		review.ProductID,
		review.Rating,
		review.Title,
		review.Comment,
		review.Status,
	).Scan(
		&review.ID,
		&review.CreatedAt,
		&review.UpdatedAt,
	)
}

func (r *ReviewRepository) GetByID(ctx context.Context, id int64) (*models.Review, error) {
	query := `
		SELECT
			id,
			user_id,
			product_id,
			rating,
			title,
			comment,
			status,
			created_at,
			updated_at
		FROM reviews
		WHERE id = $1
	`

	review := &models.Review{}

	err := r.db.QueryRowContext(
		ctx,
		query,
		id,
	).Scan(
		&review.ID,
		&review.UserID,
		&review.ProductID,
		&review.Rating,
		&review.Title,
		&review.Comment,
		&review.Status,
		&review.CreatedAt,
		&review.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return review, nil
}

func (r *ReviewRepository) GetByUserAndProduct(ctx context.Context, userID int64, productID int64) (*models.Review, error) {
	query := `
		SELECT
			id,
			user_id,
			product_id,
			rating,
			title,
			comment,
			status,
			created_at,
			updated_at
		FROM reviews
		WHERE user_id = $1
		  AND product_id = $2
	`

	review := &models.Review{}

	err := r.db.QueryRowContext(
		ctx,
		query,
		userID,
		productID,
	).Scan(
		&review.ID,
		&review.UserID,
		&review.ProductID,
		&review.Rating,
		&review.Title,
		&review.Comment,
		&review.Status,
		&review.CreatedAt,
		&review.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return review, nil
}

func (r *ReviewRepository) GetAllByProductID(ctx context.Context, productID int64) ([]*models.Review, error) {
	query := `
		SELECT
			id,
			user_id,
			product_id,
			rating,
			title,
			comment,
			status,
			created_at,
			updated_at
		FROM reviews
		WHERE product_id = $1
		  AND status = 'approved'
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reviews := make([]*models.Review, 0)

	for rows.Next() {
		review := &models.Review{}

		if err := rows.Scan(
			&review.ID,
			&review.UserID,
			&review.ProductID,
			&review.Rating,
			&review.Title,
			&review.Comment,
			&review.Status,
			&review.CreatedAt,
			&review.UpdatedAt,
		); err != nil {
			return nil, err
		}

		reviews = append(reviews, review)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return reviews, nil
}

func (r *ReviewRepository) GetAllByUserID(ctx context.Context, userID int64) ([]*models.Review, error) {
	query := `
		SELECT
			id,
			user_id,
			product_id,
			rating,
			title,
			comment,
			status,
			created_at,
			updated_at
		FROM reviews
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reviews := make([]*models.Review, 0)

	for rows.Next() {
		review := &models.Review{}

		if err := rows.Scan(
			&review.ID,
			&review.UserID,
			&review.ProductID,
			&review.Rating,
			&review.Title,
			&review.Comment,
			&review.Status,
			&review.CreatedAt,
			&review.UpdatedAt,
		); err != nil {
			return nil, err
		}

		reviews = append(reviews, review)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return reviews, nil
}

func (r *ReviewRepository) Update(ctx context.Context, review *models.Review) error {
	query := `
		UPDATE reviews
		SET
			rating = $1,
			title = $2,
			comment = $3,
			updated_at = NOW()
		WHERE id = $4
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		review.Rating,
		review.Title,
		review.Comment,
		review.ID,
	)

	return err
}

func (r *ReviewRepository) Delete(ctx context.Context, id int64) error {
	query := `
		DELETE FROM reviews
		WHERE id = $1
	`

	_, err := r.db.ExecContext(ctx, query, id)

	return err
}

func (r *ReviewRepository) GetAll(ctx context.Context) ([]*models.Review, error) {
	query := `
		SELECT
			id,
			user_id,
			product_id,
			rating,
			title,
			comment,
			status,
			created_at,
			updated_at
		FROM reviews
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reviews := make([]*models.Review, 0)

	for rows.Next() {
		review := &models.Review{}

		if err := rows.Scan(
			&review.ID,
			&review.UserID,
			&review.ProductID,
			&review.Rating,
			&review.Title,
			&review.Comment,
			&review.Status,
			&review.CreatedAt,
			&review.UpdatedAt,
		); err != nil {
			return nil, err
		}

		reviews = append(reviews, review)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return reviews, nil
}

func (r *ReviewRepository) GetAllByStatus(ctx context.Context, status string) ([]*models.Review, error) {
	query := `
		SELECT
			id,
			user_id,
			product_id,
			rating,
			title,
			comment,
			status,
			created_at,
			updated_at
		FROM reviews
		WHERE status = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reviews := make([]*models.Review, 0)

	for rows.Next() {
		review := &models.Review{}

		if err := rows.Scan(
			&review.ID,
			&review.UserID,
			&review.ProductID,
			&review.Rating,
			&review.Title,
			&review.Comment,
			&review.Status,
			&review.CreatedAt,
			&review.UpdatedAt,
		); err != nil {
			return nil, err
		}

		reviews = append(reviews, review)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return reviews, nil
}

func (r *ReviewRepository) UpdateStatus(ctx context.Context, id int64, status string) error {
	query := `
		UPDATE reviews
		SET
			status = $1,
			updated_at = NOW()
		WHERE id = $2
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		status,
		id,
	)

	return err
}
