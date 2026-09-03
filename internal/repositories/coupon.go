package repositories

import (
	"context"
	"database/sql"

	"github.com/ariiiiph/ecommerce/internal/models"
)

type CouponRepository struct {
	db *sql.DB
}

func NewCouponRepository(db *sql.DB) *CouponRepository {
	return &CouponRepository{
		db: db,
	}
}

func (r *CouponRepository) Create(ctx context.Context, coupon *models.Coupon) error {
	query := `
		INSERT INTO coupons (
			code,
			discount_type,
			discount_value,
			minimum_order_amount,
			maximum_discount,
			usage_limit,
			starts_at,
			expires_at,
			is_active
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING
			id,
			code,
			discount_type,
			discount_value,
			minimum_order_amount,
			maximum_discount,
			usage_limit,
			used_count,
			starts_at,
			expires_at,
			is_active,
			created_at,
			updated_at
	`

	return r.db.QueryRowContext(
		ctx,
		query,
		coupon.Code,
		coupon.DiscountType,
		coupon.DiscountValue,
		coupon.MinimumOrderAmount,
		coupon.MaximumDiscount,
		coupon.UsageLimit,
		coupon.StartsAt,
		coupon.ExpiresAt,
		coupon.IsActive,
	).Scan(
		&coupon.ID,
		&coupon.Code,
		&coupon.DiscountType,
		&coupon.DiscountValue,
		&coupon.MinimumOrderAmount,
		&coupon.MaximumDiscount,
		&coupon.UsageLimit,
		&coupon.UsedCount,
		&coupon.StartsAt,
		&coupon.ExpiresAt,
		&coupon.IsActive,
		&coupon.CreatedAt,
		&coupon.UpdatedAt,
	)
}

func (r *CouponRepository) GetByID(ctx context.Context, id int64) (*models.Coupon, error) {
	query := `
		SELECT
			id,
			code,
			discount_type,
			discount_value,
			minimum_order_amount,
			maximum_discount,
			usage_limit,
			used_count,
			starts_at,
			expires_at,
			is_active,
			created_at,
			updated_at
		FROM coupons
		WHERE id = $1
	`

	coupon := &models.Coupon{}

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&coupon.ID,
		&coupon.Code,
		&coupon.DiscountType,
		&coupon.DiscountValue,
		&coupon.MinimumOrderAmount,
		&coupon.MaximumDiscount,
		&coupon.UsageLimit,
		&coupon.UsedCount,
		&coupon.StartsAt,
		&coupon.ExpiresAt,
		&coupon.IsActive,
		&coupon.CreatedAt,
		&coupon.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return coupon, nil
}

func (r *CouponRepository) GetByCode(ctx context.Context, code string) (*models.Coupon, error) {
	query := `
		SELECT
			id,
			code,
			discount_type,
			discount_value,
			minimum_order_amount,
			maximum_discount,
			usage_limit,
			used_count,
			starts_at,
			expires_at,
			is_active,
			created_at,
			updated_at
		FROM coupons
		WHERE code = $1
	`

	coupon := &models.Coupon{}

	err := r.db.QueryRowContext(ctx, query, code).Scan(
		&coupon.ID,
		&coupon.Code,
		&coupon.DiscountType,
		&coupon.DiscountValue,
		&coupon.MinimumOrderAmount,
		&coupon.MaximumDiscount,
		&coupon.UsageLimit,
		&coupon.UsedCount,
		&coupon.StartsAt,
		&coupon.ExpiresAt,
		&coupon.IsActive,
		&coupon.CreatedAt,
		&coupon.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return coupon, nil
}

func (r *CouponRepository) GetAll(ctx context.Context) ([]*models.Coupon, error) {
	query := `
		SELECT
			id,
			code,
			discount_type,
			discount_value,
			minimum_order_amount,
			maximum_discount,
			usage_limit,
			used_count,
			starts_at,
			expires_at,
			is_active,
			created_at,
			updated_at
		FROM coupons
		ORDER BY id ASC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var coupons []*models.Coupon

	for rows.Next() {
		coupon := &models.Coupon{}

		if err := rows.Scan(
			&coupon.ID,
			&coupon.Code,
			&coupon.DiscountType,
			&coupon.DiscountValue,
			&coupon.MinimumOrderAmount,
			&coupon.MaximumDiscount,
			&coupon.UsageLimit,
			&coupon.UsedCount,
			&coupon.StartsAt,
			&coupon.ExpiresAt,
			&coupon.IsActive,
			&coupon.CreatedAt,
			&coupon.UpdatedAt,
		); err != nil {
			return nil, err
		}

		coupons = append(coupons, coupon)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return coupons, nil
}

func (r *CouponRepository) Update(ctx context.Context, coupon *models.Coupon) error {
	query := `
		UPDATE coupons
		SET
			code = $1,
			discount_type = $2,
			discount_value = $3,
			minimum_order_amount = $4,
			maximum_discount = $5,
			usage_limit = $6,
			starts_at = $7,
			expires_at = $8,
			is_active = $9,
			updated_at = NOW()
		WHERE id = $10
		RETURNING
			id,
			code,
			discount_type,
			discount_value,
			minimum_order_amount,
			maximum_discount,
			usage_limit,
			used_count,
			starts_at,
			expires_at,
			is_active,
			created_at,
			updated_at
	`

	return r.db.QueryRowContext(
		ctx,
		query,
		coupon.Code,
		coupon.DiscountType,
		coupon.DiscountValue,
		coupon.MinimumOrderAmount,
		coupon.MaximumDiscount,
		coupon.UsageLimit,
		coupon.StartsAt,
		coupon.ExpiresAt,
		coupon.IsActive,
		coupon.ID,
	).Scan(
		&coupon.ID,
		&coupon.Code,
		&coupon.DiscountType,
		&coupon.DiscountValue,
		&coupon.MinimumOrderAmount,
		&coupon.MaximumDiscount,
		&coupon.UsageLimit,
		&coupon.UsedCount,
		&coupon.StartsAt,
		&coupon.ExpiresAt,
		&coupon.IsActive,
		&coupon.CreatedAt,
		&coupon.UpdatedAt,
	)
}

func (r *CouponRepository) Delete(ctx context.Context, id int64) error {
	query := `
		DELETE FROM coupons
		WHERE id = $1
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
