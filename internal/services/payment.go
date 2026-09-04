package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/ariiiiph/ecommerce/internal/apperror"
	"github.com/ariiiiph/ecommerce/internal/dto"
	"github.com/ariiiiph/ecommerce/internal/models"
	"github.com/ariiiiph/ecommerce/internal/repositories"
)

const paymentCurrency = "USD"

type PaymentService struct {
	db          *sql.DB
	paymentRepo *repositories.PaymentRepository
	orderRepo   *repositories.OrderRepository
}

func NewPaymentService(db *sql.DB, paymentRepo *repositories.PaymentRepository, orderRepo *repositories.OrderRepository) *PaymentService {
	return &PaymentService{
		db:          db,
		paymentRepo: paymentRepo,
		orderRepo:   orderRepo,
	}
}

func (s *PaymentService) Create(ctx context.Context, userID int64, req *dto.CreatePaymentRequest) (*dto.PaymentResponse, error) {

	if userID <= 0 {
		return nil, apperror.BadRequest(
			"INVALID_USER_ID",
			"invalid user id",
		)
	}

	if req == nil {
		return nil, apperror.BadRequest(
			"INVALID_REQUEST",
			"invalid payment request",
		)
	}

	if req.OrderID <= 0 {
		return nil, apperror.BadRequest(
			"INVALID_ORDER_ID",
			"invalid order id",
		)
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	txOrderRepo := repositories.NewOrderRepository(tx)
	txPaymentRepo := repositories.NewPaymentRepository(tx)

	order, err := txOrderRepo.GetByID(ctx, req.OrderID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NotFound(
				"ORDER_NOT_FOUND",
				"order not found",
			)
		}

		return nil, err
	}

	if order.UserID != userID {
		return nil, apperror.NotFound(
			"ORDER_NOT_FOUND",
			"order not found",
		)
	}

	if order.Status != "pending" {
		return nil, apperror.BadRequest(
			"ORDER_NOT_PAYABLE",
			"order is not payable",
		)
	}

	_, err = txPaymentRepo.GetByOrderID(ctx, order.ID)
	if err == nil {
		return nil, apperror.Conflict(
			"PAYMENT_ALREADY_EXISTS",
			"payment already exists for this order",
		)
	}

	if !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	now := time.Now()

	transactionID := fmt.Sprintf(
		"MOCK-%d",
		now.UnixNano(),
	)

	payment := &models.Payment{
		OrderID:       order.ID,
		Amount:        order.TotalAmount,
		Currency:      paymentCurrency,
		Provider:      "mock",
		TransactionID: &transactionID,
		Status:        "paid",
		PaidAt:        &now,
	}

	if err := txPaymentRepo.Create(ctx, payment); err != nil {
		return nil, err
	}

	if err := txOrderRepo.UpdateStatus(
		ctx,
		order.ID,
		"confirmed",
	); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return toPaymentResponse(payment), nil
}

func (s *PaymentService) GetByID(ctx context.Context, userID int64, id int64) (*dto.PaymentResponse, error) {

	if userID <= 0 {
		return nil, apperror.BadRequest(
			"INVALID_USER_ID",
			"invalid user id",
		)
	}

	if id <= 0 {
		return nil, apperror.BadRequest(
			"INVALID_PAYMENT_ID",
			"invalid payment id",
		)
	}

	payment, err := s.paymentRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NotFound(
				"PAYMENT_NOT_FOUND",
				"payment not found",
			)
		}

		return nil, err
	}

	order, err := s.orderRepo.GetByID(ctx, payment.OrderID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NotFound(
				"PAYMENT_NOT_FOUND",
				"payment not found",
			)
		}

		return nil, err
	}

	if order.UserID != userID {
		return nil, apperror.NotFound(
			"PAYMENT_NOT_FOUND",
			"payment not found",
		)
	}

	return toPaymentResponse(payment), nil
}

func (s *PaymentService) GetByOrderID(ctx context.Context, userID int64, orderID int64) (*dto.PaymentResponse, error) {

	if userID <= 0 {
		return nil, apperror.BadRequest(
			"INVALID_USER_ID",
			"invalid user id",
		)
	}

	if orderID <= 0 {
		return nil, apperror.BadRequest(
			"INVALID_ORDER_ID",
			"invalid order id",
		)
	}

	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NotFound(
				"ORDER_NOT_FOUND",
				"order not found",
			)
		}

		return nil, err
	}

	if order.UserID != userID {
		return nil, apperror.NotFound(
			"ORDER_NOT_FOUND",
			"order not found",
		)
	}

	payment, err := s.paymentRepo.GetByOrderID(ctx, orderID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NotFound(
				"PAYMENT_NOT_FOUND",
				"payment not found",
			)
		}

		return nil, err
	}

	return toPaymentResponse(payment), nil
}

func toPaymentResponse(payment *models.Payment) *dto.PaymentResponse {

	return &dto.PaymentResponse{
		ID:            payment.ID,
		OrderID:       payment.OrderID,
		Amount:        payment.Amount,
		Currency:      payment.Currency,
		Provider:      payment.Provider,
		TransactionID: payment.TransactionID,
		Status:        payment.Status,
		PaidAt:        payment.PaidAt,
		CreatedAt:     payment.CreatedAt,
		UpdatedAt:     payment.UpdatedAt,
	}
}
