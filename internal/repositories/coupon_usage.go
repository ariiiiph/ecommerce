package repositories

import (
	"context"
	"database/sql"

	"github.com/ariiiiph/ecommerce/internal/db"
	"github.com/ariiiiph/ecommerce/internal/models"
)

type CouponUsageRepository struct {
	db db.DBTX
}

func NewCouponUsageRepository(db db.DBTX) *CouponUsageRepository {
	return &CouponUsageRepository{
		db: db,
	}
}

func (r *CouponUsageRepository) Create(ctx context.Context, usage *models.CouponUsage) error {
	query := `
		INSERT INTO coupon_usages (
			coupon_id,
			user_id,
			order_id
		)
		VALUES ($1, $2, $3)
		RETURNING id, used_at
	`

	return r.db.QueryRowContext(
		ctx,
		query,
		usage.CouponID,
		usage.UserID,
		usage.OrderID,
	).Scan(
		&usage.ID,
		&usage.UsedAt,
	)
}

func (r *CouponUsageRepository) GetByID(ctx context.Context, id int64) (*models.CouponUsage, error) {
	query := `
		SELECT
			id,
			coupon_id,
			user_id,
			order_id,
			used_at
		FROM coupon_usages
		WHERE id = $1
	`

	usage := &models.CouponUsage{}

	err := r.db.QueryRowContext(
		ctx,
		query,
		id,
	).Scan(
		&usage.ID,
		&usage.CouponID,
		&usage.UserID,
		&usage.OrderID,
		&usage.UsedAt,
	)

	if err != nil {
		return nil, err
	}

	return usage, nil
}

func (r *CouponUsageRepository) GetByOrderID(ctx context.Context, orderID int64) (*models.CouponUsage, error) {
	query := `
		SELECT
			id,
			coupon_id,
			user_id,
			order_id,
			used_at
		FROM coupon_usages
		WHERE order_id = $1
	`

	usage := &models.CouponUsage{}

	err := r.db.QueryRowContext(
		ctx,
		query,
		orderID,
	).Scan(
		&usage.ID,
		&usage.CouponID,
		&usage.UserID,
		&usage.OrderID,
		&usage.UsedAt,
	)

	if err != nil {
		return nil, err
	}

	return usage, nil
}

func (r *CouponUsageRepository) GetAllByUserID(ctx context.Context, userID int64) ([]*models.CouponUsage, error) {
	query := `
		SELECT
			id,
			coupon_id,
			user_id,
			order_id,
			used_at
		FROM coupon_usages
		WHERE user_id = $1
		ORDER BY id DESC
	`

	rows, err := r.db.QueryContext(
		ctx,
		query,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	usages := make([]*models.CouponUsage, 0)

	for rows.Next() {
		usage := &models.CouponUsage{}

		err := rows.Scan(
			&usage.ID,
			&usage.CouponID,
			&usage.UserID,
			&usage.OrderID,
			&usage.UsedAt,
		)

		if err != nil {
			return nil, err
		}

		usages = append(usages, usage)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return usages, nil
}

func (r *CouponUsageRepository) Delete(ctx context.Context, id int64) error {
	query := `
		DELETE FROM coupon_usages
		WHERE id = $1
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		id,
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
