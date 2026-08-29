package services

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ariiiiph/ecommerce/internal/apperror"
	"github.com/ariiiiph/ecommerce/internal/dto"
	"github.com/ariiiiph/ecommerce/internal/models"
	"github.com/ariiiiph/ecommerce/internal/repositories"
)

type CartService struct {
	cartRepo *repositories.CartRepository
	userRepo *repositories.UserRepository
}

func NewCartService(cartRepo *repositories.CartRepository, userRepo *repositories.UserRepository) *CartService {
	return &CartService{
		cartRepo: cartRepo,
		userRepo: userRepo,
	}
}

func (s *CartService) Create(ctx context.Context, userID int64) (*dto.CartResponse, error) {
	if userID <= 0 {
		return nil, apperror.BadRequest(
			"INVALID_USER_ID",
			"invalid user id",
		)
	}
	_, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NotFound(
				"USER_NOT_FOUND",
				"user not found",
			)
		}
		return nil, err
	}

	_, err = s.cartRepo.GetByUserID(ctx, userID)
	if err == nil {
		return nil, apperror.Conflict(
			"CART_ALREADY_EXISTS",
			"cart already exists",
		)
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	cart := &models.Cart{
		UserID: userID,
	}

	if err := s.cartRepo.Create(ctx, cart); err != nil {
		return nil, err
	}

	return toCartResponse(cart), nil
}

func (s *CartService) GetByID(ctx context.Context, id int64, userID int64) (*dto.CartResponse, error) {
	if id <= 0 {
		return nil, apperror.BadRequest(
			"INVALID_CART_ID",
			"invalid cart id",
		)
	}

	cart, err := s.cartRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NotFound(
				"CART_NOT_FOUND",
				"cart not found",
			)
		}
		return nil, err
	}

	if cart.UserID != userID {
		return nil, apperror.Forbidden(
			"FORBIDDEN",
			"you do not have access to this cart",
		)
	}

	return toCartResponse(cart), nil
}

func (s *CartService) GetByUserID(ctx context.Context, userID int64) (*dto.CartResponse, error) {
	if userID <= 0 {
		return nil, apperror.BadRequest(
			"INVALID_USER_ID",
			"invalid user id",
		)
	}

	cart, err := s.cartRepo.GetByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NotFound(
				"CART_NOT_FOUND",
				"cart not found",
			)
		}
		return nil, err
	}
	return toCartResponse(cart), nil
}

func (s *CartService) Delete(ctx context.Context, id int64, userID int64) error {
	if id <= 0 {
		return apperror.BadRequest(
			"INVALID_CART_ID",
			"invalid cart id",
		)
	}

	if userID <= 0 {
		return apperror.BadRequest(
			"INVALID_USER_ID",
			"invalid user id",
		)
	}

	cart, err := s.cartRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return apperror.NotFound(
				"CART_NOT_FOUND",
				"cart not found",
			)
		}

		return err
	}

	if cart.UserID != userID {
		return apperror.Forbidden(
			"FORBIDDEN",
			"you do not have access to this cart",
		)
	}

	return s.cartRepo.Delete(ctx, id)
}

func toCartResponse(cart *models.Cart) *dto.CartResponse {
	return &dto.CartResponse{
		ID:        cart.ID,
		UserID:    cart.UserID,
		CreatedAt: cart.CreatedAt,
		UpdatedAt: cart.UpdatedAt,
	}
}
