package repositories

import (
	"context"
	"time"

	"github.com/ariiiiph/ecommerce/internal/db"
	"github.com/ariiiiph/ecommerce/internal/models"
)

type PaymentRepository struct {
	db db.DBTX
}

func NewPaymentRepository(db db.DBTX) *PaymentRepository {
	return &PaymentRepository{
		db: db,
	}
}

func (r *PaymentRepository) Create(ctx context.Context, payment *models.Payment) error {

	query := `
		INSERT INTO payments (
			order_id,
			amount,
			currency,
			provider,
			transaction_id,
			status,
			paid_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
	`

	return r.db.QueryRowContext(
		ctx,
		query,
		payment.OrderID,
		payment.Amount,
		payment.Currency,
		payment.Provider,
		payment.TransactionID,
		payment.Status,
		payment.PaidAt,
	).Scan(
		&payment.ID,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)
}

func (r *PaymentRepository) GetByID(ctx context.Context, id int64) (*models.Payment, error) {

	query := `
		SELECT
			id,
			order_id,
			amount,
			currency,
			provider,
			transaction_id,
			status,
			paid_at,
			created_at,
			updated_at
		FROM payments
		WHERE id = $1
	`

	payment := &models.Payment{}

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&payment.ID,
		&payment.OrderID,
		&payment.Amount,
		&payment.Currency,
		&payment.Provider,
		&payment.TransactionID,
		&payment.Status,
		&payment.PaidAt,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return payment, nil
}

func (r *PaymentRepository) GetByOrderID(ctx context.Context, orderID int64) (*models.Payment, error) {

	query := `
		SELECT
			id,
			order_id,
			amount,
			currency,
			provider,
			transaction_id,
			status,
			paid_at,
			created_at,
			updated_at
		FROM payments
		WHERE order_id = $1
	`

	payment := &models.Payment{}

	err := r.db.QueryRowContext(ctx, query, orderID).Scan(
		&payment.ID,
		&payment.OrderID,
		&payment.Amount,
		&payment.Currency,
		&payment.Provider,
		&payment.TransactionID,
		&payment.Status,
		&payment.PaidAt,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return payment, nil
}

func (r *PaymentRepository) UpdateStatus(ctx context.Context, id int64, status string, paidAt *time.Time) error {

	query := `
		UPDATE payments
		SET
			status = $1,
			paid_at = $2,
			updated_at = NOW()
		WHERE id = $3
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		status,
		paidAt,
		id,
	)

	return err
}
