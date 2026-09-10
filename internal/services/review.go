package services

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/ariiiiph/ecommerce/internal/apperror"
	"github.com/ariiiiph/ecommerce/internal/dto"
	"github.com/ariiiiph/ecommerce/internal/models"
	"github.com/ariiiiph/ecommerce/internal/repositories"
)

type ReviewService struct {
	reviewRepo  *repositories.ReviewRepository
	productRepo *repositories.ProductRepository
}

func NewReviewService(reviewRepo *repositories.ReviewRepository, productRepo *repositories.ProductRepository) *ReviewService {
	return &ReviewService{
		reviewRepo:  reviewRepo,
		productRepo: productRepo,
	}
}

func (s *ReviewService) Create(ctx context.Context, userID int64, req *dto.CreateReviewRequest) (*dto.ReviewResponse, error) {

	if userID <= 0 {
		return nil, apperror.BadRequest(
			"INVALID_USER_ID",
			"invalid user id",
		)
	}

	if req == nil {
		return nil, apperror.BadRequest(
			"INVALID_REQUEST",
			"invalid review request",
		)
	}

	if req.ProductID <= 0 {
		return nil, apperror.BadRequest(
			"INVALID_PRODUCT_ID",
			"invalid product id",
		)
	}

	if req.Rating < 1 || req.Rating > 5 {
		return nil, apperror.BadRequest(
			"INVALID_REVIEW_RATING",
			"rating must be between 1 and 5",
		)
	}

	if _, err := s.productRepo.GetByID(ctx, req.ProductID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NotFound(
				"PRODUCT_NOT_FOUND",
				"product not found",
			)
		}

		return nil, err
	}

	_, err := s.reviewRepo.GetByUserAndProduct(
		ctx,
		userID,
		req.ProductID,
	)

	if err == nil {
		return nil, apperror.Conflict(
			"REVIEW_ALREADY_EXISTS",
			"review already exists for this product",
		)
	}

	if !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	var title *string
	if req.Title != nil {
		value := strings.TrimSpace(*req.Title)
		if value != "" {
			title = &value
		}
	}

	var comment *string
	if req.Comment != nil {
		value := strings.TrimSpace(*req.Comment)
		if value != "" {
			comment = &value
		}
	}

	review := &models.Review{
		UserID:    userID,
		ProductID: req.ProductID,
		Rating:    req.Rating,
		Title:     title,
		Comment:   comment,
		Status:    "pending",
	}

	if err := s.reviewRepo.Create(ctx, review); err != nil {
		return nil, err
	}

	return toReviewResponse(review), nil
}

func (s *ReviewService) GetByID(ctx context.Context, id int64, userID int64) (*dto.ReviewResponse, error) {

	if id <= 0 {
		return nil, apperror.BadRequest(
			"INVALID_REVIEW_ID",
			"invalid review id",
		)
	}

	if userID <= 0 {
		return nil, apperror.BadRequest(
			"INVALID_USER_ID",
			"invalid user id",
		)
	}

	review, err := s.reviewRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NotFound(
				"REVIEW_NOT_FOUND",
				"review not found",
			)
		}

		return nil, err
	}

	if review.UserID != userID {
		return nil, apperror.NotFound(
			"REVIEW_NOT_FOUND",
			"review not found",
		)
	}

	return toReviewResponse(review), nil
}

func (s *ReviewService) GetAllByProductID(ctx context.Context, productID int64) ([]*dto.ReviewResponse, error) {

	if productID <= 0 {
		return nil, apperror.BadRequest(
			"INVALID_PRODUCT_ID",
			"invalid product id",
		)
	}

	reviews, err := s.reviewRepo.GetAllByProductID(ctx, productID)
	if err != nil {
		return nil, err
	}

	result := make([]*dto.ReviewResponse, 0, len(reviews))

	for _, review := range reviews {
		result = append(result, toReviewResponse(review))
	}

	return result, nil
}

func (s *ReviewService) GetAllByUserID(ctx context.Context, userID int64) ([]*dto.ReviewResponse, error) {

	if userID <= 0 {
		return nil, apperror.BadRequest(
			"INVALID_USER_ID",
			"invalid user id",
		)
	}

	reviews, err := s.reviewRepo.GetAllByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	result := make([]*dto.ReviewResponse, 0, len(reviews))

	for _, review := range reviews {
		result = append(result, toReviewResponse(review))
	}

	return result, nil
}

func (s *ReviewService) Update(ctx context.Context, id int64, userID int64, req *dto.UpdateReviewRequest) (*dto.ReviewResponse, error) {

	if id <= 0 {
		return nil, apperror.BadRequest(
			"INVALID_REVIEW_ID",
			"invalid review id",
		)
	}

	if userID <= 0 {
		return nil, apperror.BadRequest(
			"INVALID_USER_ID",
			"invalid user id",
		)
	}

	if req == nil {
		return nil, apperror.BadRequest(
			"INVALID_REQUEST",
			"invalid review request",
		)
	}

	review, err := s.reviewRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NotFound(
				"REVIEW_NOT_FOUND",
				"review not found",
			)
		}

		return nil, err
	}

	if review.UserID != userID {
		return nil, apperror.NotFound(
			"REVIEW_NOT_FOUND",
			"review not found",
		)
	}

	if req.Rating != nil {
		if *req.Rating < 1 || *req.Rating > 5 {
			return nil, apperror.BadRequest(
				"INVALID_REVIEW_RATING",
				"rating must be between 1 and 5",
			)
		}

		review.Rating = *req.Rating
	}

	if req.Title != nil {
		value := strings.TrimSpace(*req.Title)

		if value == "" {
			review.Title = nil
		} else {
			review.Title = &value
		}
	}

	if req.Comment != nil {
		value := strings.TrimSpace(*req.Comment)

		if value == "" {
			review.Comment = nil
		} else {
			review.Comment = &value
		}
	}

	review.Status = "pending"

	if err := s.reviewRepo.Update(ctx, review); err != nil {
		return nil, err
	}

	return toReviewResponse(review), nil
}

func (s *ReviewService) Delete(ctx context.Context, id int64, userID int64) error {

	if id <= 0 {
		return apperror.BadRequest(
			"INVALID_REVIEW_ID",
			"invalid review id",
		)
	}

	if userID <= 0 {
		return apperror.BadRequest(
			"INVALID_USER_ID",
			"invalid user id",
		)
	}

	review, err := s.reviewRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return apperror.NotFound(
				"REVIEW_NOT_FOUND",
				"review not found",
			)
		}

		return err
	}

	if review.UserID != userID {
		return apperror.NotFound(
			"REVIEW_NOT_FOUND",
			"review not found",
		)
	}

	return s.reviewRepo.Delete(ctx, id)
}

func (s *ReviewService) GetAllForAdmin(ctx context.Context) ([]*dto.ReviewResponse, error) {

	reviews, err := s.reviewRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]*dto.ReviewResponse, 0, len(reviews))

	for _, review := range reviews {
		result = append(result, toReviewResponse(review))
	}

	return result, nil
}

func (s *ReviewService) GetAllByStatusForAdmin(ctx context.Context, status string) ([]*dto.ReviewResponse, error) {

	status = strings.TrimSpace(status)

	if status != "pending" &&
		status != "approved" &&
		status != "rejected" {
		return nil, apperror.BadRequest(
			"INVALID_REVIEW_STATUS",
			"invalid review status",
		)
	}

	reviews, err := s.reviewRepo.GetAllByStatus(ctx, status)
	if err != nil {
		return nil, err
	}

	result := make([]*dto.ReviewResponse, 0, len(reviews))

	for _, review := range reviews {
		result = append(result, toReviewResponse(review))
	}

	return result, nil
}

func (s *ReviewService) Approve(ctx context.Context, id int64) error {

	if id <= 0 {
		return apperror.BadRequest(
			"INVALID_REVIEW_ID",
			"invalid review id",
		)
	}

	review, err := s.reviewRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return apperror.NotFound(
				"REVIEW_NOT_FOUND",
				"review not found",
			)
		}

		return err
	}

	if review.Status == "approved" {
		return apperror.Conflict(
			"REVIEW_ALREADY_APPROVED",
			"review is already approved",
		)
	}

	if review.Status == "rejected" {
		return apperror.Conflict(
			"REVIEW_ALREADY_REJECTED",
			"rejected review cannot be approved",
		)
	}

	return s.reviewRepo.UpdateStatus(ctx, id, "approved")
}

func (s *ReviewService) Reject(ctx context.Context, id int64) error {

	if id <= 0 {
		return apperror.BadRequest(
			"INVALID_REVIEW_ID",
			"invalid review id",
		)
	}

	review, err := s.reviewRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return apperror.NotFound(
				"REVIEW_NOT_FOUND",
				"review not found",
			)
		}

		return err
	}

	if review.Status == "rejected" {
		return apperror.Conflict(
			"REVIEW_ALREADY_REJECTED",
			"review is already rejected",
		)
	}

	if review.Status == "approved" {
		return apperror.Conflict(
			"REVIEW_ALREADY_APPROVED",
			"approved review cannot be rejected",
		)
	}

	return s.reviewRepo.UpdateStatus(ctx, id, "rejected")
}

func toReviewResponse(review *models.Review) *dto.ReviewResponse {
	return &dto.ReviewResponse{
		ID:        review.ID,
		UserID:    review.UserID,
		ProductID: review.ProductID,
		Rating:    review.Rating,
		Title:     review.Title,
		Comment:   review.Comment,
		Status:    review.Status,
		CreatedAt: review.CreatedAt,
		UpdatedAt: review.UpdatedAt,
	}
}
