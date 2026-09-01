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

type CartItemService struct {
	cartItemRepo       *repositories.CartItemRepository
	cartRepo           *repositories.CartRepository
	productVariantRepo *repositories.ProductVariantRepository
}

func NewCartItemService(cartItemRepo *repositories.CartItemRepository, cartRepo *repositories.CartRepository, productVariantRepo *repositories.ProductVariantRepository) *CartItemService {
	return &CartItemService{
		cartItemRepo:       cartItemRepo,
		cartRepo:           cartRepo,
		productVariantRepo: productVariantRepo,
	}
}

func (s *CartItemService) Create(ctx context.Context, req *dto.CreateCartItemRequest, userID int64) (*dto.CartItemResponse, error) {
	if userID <= 0 {
		return nil, apperror.BadRequest(
			"INVALID_USER_ID",
			"invalid user id",
		)
	}

	if req.VariantID <= 0 {
		return nil, apperror.BadRequest(
			"INVALID_VARIANT_ID",
			"invalid variant id",
		)
	}

	if req.Quantity <= 0 {
		return nil, apperror.BadRequest(
			"INVALID_QUANTITY",
			"quantity must be greater than zero",
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

	_, err = s.productVariantRepo.GetByID(ctx, req.VariantID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NotFound(
				"VARIANT_NOT_FOUND",
				"variant not found",
			)
		}
		return nil, err
	}

	cartItem := &models.CartItem{
		CartID:    cart.ID,
		VariantID: req.VariantID,
		Quantity:  req.Quantity,
	}

	if err := s.cartItemRepo.Create(ctx, cartItem); err != nil {
		return nil, err
	}

	return toCartItemResponse(cartItem), nil
}

func (s *CartItemService) GetByID(ctx context.Context, id int64, userID int64) (*dto.CartItemResponse, error) {
	if id <= 0 {
		return nil, apperror.BadRequest(
			"INVALID_CART_ITEM_ID",
			"invalid cart item id",
		)
	}

	if userID <= 0 {
		return nil, apperror.BadRequest(
			"INVALID_USER_ID",
			"invalid user id",
		)
	}

	cartItem, err := s.cartItemRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NotFound(
				"CART_ITEM_NOT_FOUND",
				"cart item not found",
			)
		}
		return nil, err
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

	if cartItem.CartID != cart.ID {
		return nil, apperror.Forbidden(
			"FORBIDDEN",
			"you do not have access to this cart item",
		)
	}

	return toCartItemResponse(cartItem), nil
}

func (s *CartItemService) GetAllByCartID(ctx context.Context, cartID int64, userID int64) ([]*dto.CartItemResponse, error) {
	if cartID <= 0 {
		return nil, apperror.BadRequest(
			"INVALID_CART_ID",
			"invalid cart id",
		)
	}

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

	if cart.ID != cartID {
		return nil, apperror.Forbidden(
			"FORBIDDEN",
			"you do not have access to this cart",
		)
	}

	cartItems, err := s.cartItemRepo.GetAllByCartID(ctx, cartID)
	if err != nil {
		return nil, err
	}

	result := make([]*dto.CartItemResponse, 0, len(cartItems))

	for _, cartItem := range cartItems {
		result = append(result, toCartItemResponse(cartItem))
	}

	return result, nil
}

func (s *CartItemService) Update(ctx context.Context, id int64, req *dto.UpdateCartItemRequest, userID int64) (*dto.CartItemResponse, error) {
	if id <= 0 {
		return nil, apperror.BadRequest(
			"INVALID_CART_ITEM_ID",
			"invalid cart item id",
		)
	}

	if userID <= 0 {
		return nil, apperror.BadRequest(
			"INVALID_USER_ID",
			"invalid user id",
		)
	}

	if req.Quantity <= 0 {
		return nil, apperror.BadRequest(
			"INVALID_QUANTITY",
			"quantity must be greater than zero",
		)
	}

	cartItem, err := s.cartItemRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NotFound(
				"CART_ITEM_NOT_FOUND",
				"cart item not found",
			)
		}
		return nil, err
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

	if cartItem.CartID != cart.ID {
		return nil, apperror.Forbidden(
			"FORBIDDEN",
			"you do not have access to this cart item",
		)
	}

	cartItem.Quantity = req.Quantity

	if err := s.cartItemRepo.Update(ctx, cartItem); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NotFound(
				"CART_ITEM_NOT_FOUND",
				"cart item not found",
			)
		}
		return nil, err
	}

	return toCartItemResponse(cartItem), nil
}

func (s *CartItemService) Delete(ctx context.Context, id int64, userID int64) error {
	if id <= 0 {
		return apperror.BadRequest(
			"INVALID_CART_ITEM_ID",
			"invalid cart item id",
		)
	}

	if userID <= 0 {
		return apperror.BadRequest(
			"INVALID_USER_ID",
			"invalid user id",
		)
	}

	cartItem, err := s.cartItemRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return apperror.NotFound(
				"CART_ITEM_NOT_FOUND",
				"cart item not found",
			)
		}
		return err
	}

	cart, err := s.cartRepo.GetByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return apperror.NotFound(
				"CART_NOT_FOUND",
				"cart not found",
			)
		}
		return err
	}

	if cartItem.CartID != cart.ID {
		return apperror.Forbidden(
			"FORBIDDEN",
			"you do not have access to this cart item",
		)
	}

	if err := s.cartItemRepo.Delete(ctx, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return apperror.NotFound(
				"CART_ITEM_NOT_FOUND",
				"cart item not found",
			)
		}
		return err
	}

	return nil
}

func toCartItemResponse(cartItem *models.CartItem) *dto.CartItemResponse {
	return &dto.CartItemResponse{
		ID:        cartItem.ID,
		CartID:    cartItem.CartID,
		VariantID: cartItem.VariantID,
		Quantity:  cartItem.Quantity,
		CreatedAt: cartItem.CreatedAt,
		UpdatedAt: cartItem.UpdatedAt,
	}
}
