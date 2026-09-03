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

type WishlistService struct {
	wishlistRepo *repositories.WishlistRepository
	userRepo     *repositories.UserRepository
}

func NewWishlistService(wishlistRepo *repositories.WishlistRepository, userRepo *repositories.UserRepository) *WishlistService {
	return &WishlistService{
		wishlistRepo: wishlistRepo,
		userRepo:     userRepo,
	}
}
func (s *WishlistService) Create(ctx context.Context, userID int64) (*dto.WishlistResponse, error) {
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
	_, err = s.wishlistRepo.GetByUserID(ctx, userID)
	if err == nil {
		return nil, apperror.Conflict(
			"WISHLIST_ALREADY_EXISTS",
			"wishlist already exists",
		)
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	wishlist := &models.Wishlist{
		UserID: userID,
	}

	if err := s.wishlistRepo.Create(ctx, wishlist); err != nil {
		return nil, err
	}

	return toWishlistResponse(wishlist), nil
}

func (s *WishlistService) GetByID(ctx context.Context, id, userID int64) (*dto.WishlistResponse, error) {
	if id <= 0 {
		return nil, apperror.BadRequest(
			"INVALID_WISHLIST_ID",
			"invalid wishlist id",
		)
	}

	wishlist, err := s.wishlistRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NotFound(
				"WISHLIST_NOT_FOUND",
				"wishlist not found",
			)
		}
		return nil, err
	}

	if wishlist.UserID != userID {
		return nil, apperror.Forbidden(
			"FORBIDDEN",
			"you do not have access to this wishlist",
		)
	}

	return toWishlistResponse(wishlist), nil

}

func (s *WishlistService) GetByUserID(ctx context.Context, userID int64) (*dto.WishlistResponse, error) {
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
	return toWishlistResponse(wishlist), nil
}

func (s *WishlistService) Delete(ctx context.Context, id, userID int64) error {
	if id <= 0 {
		return apperror.BadRequest(
			"INVALID_WISHLIST_ID",
			"invalid wishlist id",
		)
	}

	if userID <= 0 {
		return apperror.BadRequest(
			"INVALID_USER_ID",
			"invalid user id",
		)
	}

	wishlist, err := s.wishlistRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return apperror.NotFound(
				"WISHLIST_NOT_FOUND",
				"wishlist not found",
			)
		}
		return err
	}

	if wishlist.UserID != userID {
		return apperror.Forbidden(
			"FORBIDDEN",
			"you do not have access to this wishlist",
		)
	}

	return s.wishlistRepo.Delete(ctx, id)

}

func toWishlistResponse(wishlist *models.Wishlist) *dto.WishlistResponse {
	return &dto.WishlistResponse{
		ID:        wishlist.ID,
		UserID:    wishlist.UserID,
		CreatedAt: wishlist.CreatedAt,
	}
}
