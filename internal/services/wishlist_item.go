package services

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ariiiiph/ecommerce/internal/apperror"
	"github.com/ariiiiph/ecommerce/internal/dto"
	"github.com/ariiiiph/ecommerce/internal/models"
	"github.com/ariiiiph/ecommerce/internal/repositories"
	"github.com/lib/pq"
)

type WishlistItemService struct {
	wishlistItemRepo *repositories.WishlistItemRepository
	wishlistRepo     *repositories.WishlistRepository
	productRepo      *repositories.ProductRepository
}

func NewWishlistItemService(wishlistItemRepo *repositories.WishlistItemRepository, wishlistRepo *repositories.WishlistRepository, productRepo *repositories.ProductRepository) *WishlistItemService {
	return &WishlistItemService{
		wishlistItemRepo: wishlistItemRepo,
		wishlistRepo:     wishlistRepo,
		productRepo:      productRepo,
	}
}

func (s *WishlistItemService) Create(ctx context.Context, userID int64, req *dto.CreateWishlistItemRequest) (*dto.WishlistItemResponse, error) {

	if userID <= 0 {
		return nil, apperror.BadRequest(
			"INVALID_USER_ID",
			"invalid user id",
		)
	}

	if req.ProductID <= 0 {
		return nil, apperror.BadRequest(
			"INVALID_PRODUCT_ID",
			"invalid product id",
		)
	}

	_, err := s.productRepo.GetByID(ctx, req.ProductID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NotFound(
				"PRODUCT_NOT_FOUND",
				"product not found",
			)
		}

		return nil, err
	}

	wishlist, err := s.wishlistRepo.GetByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NotFound(
				"WISHLIST_NOT_FOUND",
				"wishlist not found",
			)
		}

		return nil, err
	}

	item := &models.WishlistItem{
		WishlistID: wishlist.ID,
		ProductID:  req.ProductID,
	}

	if err := s.wishlistItemRepo.Create(ctx, item); err != nil {
		if isUniqueViolation(err) {
			return nil, apperror.Conflict(
				"WISHLIST_ITEM_ALREADY_EXISTS",
				"product already exists in wishlist",
			)
		}

		return nil, err
	}

	return toWishlistItemResponse(item), nil
}

func (s *WishlistItemService) GetByID(ctx context.Context, id int64, userID int64) (*dto.WishlistItemResponse, error) {

	if id <= 0 {
		return nil, apperror.BadRequest(
			"INVALID_WISHLIST_ITEM_ID",
			"invalid wishlist item id",
		)
	}

	if userID <= 0 {
		return nil, apperror.BadRequest(
			"INVALID_USER_ID",
			"invalid user id",
		)
	}

	item, err := s.wishlistItemRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NotFound(
				"WISHLIST_ITEM_NOT_FOUND",
				"wishlist item not found",
			)
		}

		return nil, err
	}

	wishlist, err := s.wishlistRepo.GetByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NotFound(
				"WISHLIST_NOT_FOUND",
				"wishlist not found",
			)
		}

		return nil, err
	}

	if item.WishlistID != wishlist.ID {
		return nil, apperror.Forbidden(
			"FORBIDDEN",
			"you do not have access to this wishlist item",
		)
	}

	return toWishlistItemResponse(item), nil
}

func (s *WishlistItemService) GetAllByWishlistID(ctx context.Context, wishlistID int64, userID int64) ([]*dto.WishlistItemResponse, error) {

	if wishlistID <= 0 {
		return nil, apperror.BadRequest(
			"INVALID_WISHLIST_ID",
			"invalid wishlist id",
		)
	}

	if userID <= 0 {
		return nil, apperror.BadRequest(
			"INVALID_USER_ID",
			"invalid user id",
		)
	}

	wishlist, err := s.wishlistRepo.GetByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NotFound(
				"WISHLIST_NOT_FOUND",
				"wishlist not found",
			)
		}

		return nil, err
	}

	if wishlist.ID != wishlistID {
		return nil, apperror.Forbidden(
			"FORBIDDEN",
			"you do not have access to this wishlist",
		)
	}

	items, err := s.wishlistItemRepo.GetAllByWishlistID(ctx, wishlistID)
	if err != nil {
		return nil, err
	}

	result := make([]*dto.WishlistItemResponse, 0, len(items))

	for _, item := range items {
		result = append(result, toWishlistItemResponse(item))
	}

	return result, nil
}

func (s *WishlistItemService) Delete(ctx context.Context, id int64, userID int64) error {

	if id <= 0 {
		return apperror.BadRequest(
			"INVALID_WISHLIST_ITEM_ID",
			"invalid wishlist item id",
		)
	}

	if userID <= 0 {
		return apperror.BadRequest(
			"INVALID_USER_ID",
			"invalid user id",
		)
	}

	item, err := s.wishlistItemRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return apperror.NotFound(
				"WISHLIST_ITEM_NOT_FOUND",
				"wishlist item not found",
			)
		}

		return err
	}

	wishlist, err := s.wishlistRepo.GetByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return apperror.NotFound(
				"WISHLIST_NOT_FOUND",
				"wishlist not found",
			)
		}

		return err
	}

	if item.WishlistID != wishlist.ID {
		return apperror.Forbidden(
			"FORBIDDEN",
			"you do not have access to this wishlist item",
		)
	}

	return s.wishlistItemRepo.Delete(ctx, id)
}

func toWishlistItemResponse(item *models.WishlistItem) *dto.WishlistItemResponse {
	return &dto.WishlistItemResponse{
		ID:         item.ID,
		WishlistID: item.WishlistID,
		ProductID:  item.ProductID,
		CreatedAt:  item.CreatedAt,
	}
}

func isUniqueViolation(err error) bool {
	var pqErr *pq.Error

	if errors.As(err, &pqErr) {
		return pqErr.Code == "23505"
	}

	return false
}
