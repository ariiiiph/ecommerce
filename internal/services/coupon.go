package services

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/ariiiiph/ecommerce/internal/apperror"
	"github.com/ariiiiph/ecommerce/internal/dto"
	"github.com/ariiiiph/ecommerce/internal/models"
	"github.com/ariiiiph/ecommerce/internal/repositories"
	"github.com/shopspring/decimal"
)

type CouponService struct {
	couponRepo *repositories.CouponRepository
}

func NewCouponService(couponRepo *repositories.CouponRepository) *CouponService {
	return &CouponService{
		couponRepo: couponRepo,
	}
}

func (s *CouponService) Create(ctx context.Context, req *dto.CreateCouponRequest) (*dto.CouponResponse, error) {

	code := strings.TrimSpace(req.Code)

	if code == "" {
		return nil, apperror.BadRequest(
			"COUPON_CODE_REQUIRED",
			"coupon code is required",
		)
	}

	if len(code) > 50 {
		return nil, apperror.BadRequest(
			"INVALID_COUPON_CODE",
			"coupon code must not exceed 50 characters",
		)
	}

	discountType := strings.TrimSpace(req.DiscountType)

	if discountType != "percentage" && discountType != "fixed" {
		return nil, apperror.BadRequest(
			"INVALID_DISCOUNT_TYPE",
			"discount type must be percentage or fixed",
		)
	}

	if !req.DiscountValue.GreaterThan(decimal.Zero) {
		return nil, apperror.BadRequest(
			"INVALID_DISCOUNT_VALUE",
			"discount value must be greater than zero",
		)
	}

	if discountType == "percentage" &&
		req.DiscountValue.GreaterThan(decimal.NewFromInt(100)) {
		return nil, apperror.BadRequest(
			"INVALID_DISCOUNT_VALUE",
			"percentage discount cannot exceed 100",
		)
	}

	if req.MinimumOrderAmount != nil &&
		req.MinimumOrderAmount.IsNegative() {
		return nil, apperror.BadRequest(
			"INVALID_MINIMUM_ORDER_AMOUNT",
			"minimum order amount cannot be negative",
		)
	}

	if req.MaximumDiscount != nil &&
		req.MaximumDiscount.IsNegative() {
		return nil, apperror.BadRequest(
			"INVALID_MAXIMUM_DISCOUNT",
			"maximum discount cannot be negative",
		)
	}

	if req.UsageLimit != nil && *req.UsageLimit <= 0 {
		return nil, apperror.BadRequest(
			"INVALID_USAGE_LIMIT",
			"usage limit must be greater than zero",
		)
	}

	if !req.ExpiresAt.After(req.StartsAt) {
		return nil, apperror.BadRequest(
			"INVALID_COUPON_DATES",
			"expires at must be after starts at",
		)
	}

	_, err := s.couponRepo.GetByCode(ctx, code)
	if err == nil {
		return nil, apperror.Conflict(
			"COUPON_ALREADY_EXISTS",
			"coupon already exists",
		)
	}

	if !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	coupon := &models.Coupon{
		Code:               code,
		DiscountType:       discountType,
		DiscountValue:      req.DiscountValue,
		MinimumOrderAmount: req.MinimumOrderAmount,
		MaximumDiscount:    req.MaximumDiscount,
		UsageLimit:         req.UsageLimit,
		StartsAt:           req.StartsAt,
		ExpiresAt:          req.ExpiresAt,
		IsActive:           req.IsActive,
	}

	if err := s.couponRepo.Create(ctx, coupon); err != nil {
		return nil, err
	}

	return toCouponResponse(coupon), nil
}

func (s *CouponService) GetByID(ctx context.Context, id int64) (*dto.CouponResponse, error) {

	if id <= 0 {
		return nil, apperror.BadRequest(
			"INVALID_COUPON_ID",
			"invalid coupon id",
		)
	}

	coupon, err := s.couponRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NotFound(
				"COUPON_NOT_FOUND",
				"coupon not found",
			)
		}

		return nil, err
	}

	return toCouponResponse(coupon), nil
}

func (s *CouponService) GetAll(ctx context.Context) ([]*dto.CouponResponse, error) {

	coupons, err := s.couponRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]*dto.CouponResponse, 0, len(coupons))

	for _, coupon := range coupons {
		result = append(result, toCouponResponse(coupon))
	}

	return result, nil
}

func (s *CouponService) Update(ctx context.Context, id int64, req *dto.UpdateCouponRequest) (*dto.CouponResponse, error) {

	if id <= 0 {
		return nil, apperror.BadRequest(
			"INVALID_COUPON_ID",
			"invalid coupon id",
		)
	}

	coupon, err := s.couponRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NotFound(
				"COUPON_NOT_FOUND",
				"coupon not found",
			)
		}

		return nil, err
	}

	code := strings.TrimSpace(req.Code)

	if code == "" {
		return nil, apperror.BadRequest(
			"COUPON_CODE_REQUIRED",
			"coupon code is required",
		)
	}

	if len(code) > 50 {
		return nil, apperror.BadRequest(
			"INVALID_COUPON_CODE",
			"coupon code must not exceed 50 characters",
		)
	}

	discountType := strings.TrimSpace(req.DiscountType)

	if discountType != "percentage" && discountType != "fixed" {
		return nil, apperror.BadRequest(
			"INVALID_DISCOUNT_TYPE",
			"discount type must be percentage or fixed",
		)
	}

	if !req.DiscountValue.GreaterThan(decimal.Zero) {
		return nil, apperror.BadRequest(
			"INVALID_DISCOUNT_VALUE",
			"discount value must be greater than zero",
		)
	}

	if discountType == "percentage" &&
		req.DiscountValue.GreaterThan(decimal.NewFromInt(100)) {
		return nil, apperror.BadRequest(
			"INVALID_DISCOUNT_VALUE",
			"percentage discount cannot exceed 100",
		)
	}

	if req.MinimumOrderAmount != nil &&
		req.MinimumOrderAmount.IsNegative() {
		return nil, apperror.BadRequest(
			"INVALID_MINIMUM_ORDER_AMOUNT",
			"minimum order amount cannot be negative",
		)
	}

	if req.MaximumDiscount != nil &&
		req.MaximumDiscount.IsNegative() {
		return nil, apperror.BadRequest(
			"INVALID_MAXIMUM_DISCOUNT",
			"maximum discount cannot be negative",
		)
	}

	if req.UsageLimit != nil && *req.UsageLimit <= 0 {
		return nil, apperror.BadRequest(
			"INVALID_USAGE_LIMIT",
			"usage limit must be greater than zero",
		)
	}

	if req.UsageLimit != nil && *req.UsageLimit < coupon.UsedCount {
		return nil, apperror.BadRequest(
			"INVALID_USAGE_LIMIT",
			"usage limit cannot be less than current used count",
		)
	}

	if !req.ExpiresAt.After(req.StartsAt) {
		return nil, apperror.BadRequest(
			"INVALID_COUPON_DATES",
			"expires at must be after starts at",
		)
	}

	if !strings.EqualFold(code, coupon.Code) {
		existingCoupon, err := s.couponRepo.GetByCode(ctx, code)

		if err == nil && existingCoupon.ID != id {
			return nil, apperror.Conflict(
				"COUPON_ALREADY_EXISTS",
				"coupon already exists",
			)
		}

		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}

	coupon.Code = code
	coupon.DiscountType = discountType
	coupon.DiscountValue = req.DiscountValue
	coupon.MinimumOrderAmount = req.MinimumOrderAmount
	coupon.MaximumDiscount = req.MaximumDiscount
	coupon.UsageLimit = req.UsageLimit
	coupon.StartsAt = req.StartsAt
	coupon.ExpiresAt = req.ExpiresAt
	coupon.IsActive = req.IsActive

	if err := s.couponRepo.Update(ctx, coupon); err != nil {
		return nil, err
	}

	return toCouponResponse(coupon), nil
}

func (s *CouponService) Delete(ctx context.Context, id int64) error {

	if id <= 0 {
		return apperror.BadRequest(
			"INVALID_COUPON_ID",
			"invalid coupon id",
		)
	}

	_, err := s.couponRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return apperror.NotFound(
				"COUPON_NOT_FOUND",
				"coupon not found",
			)
		}

		return err
	}

	return s.couponRepo.Delete(ctx, id)
}

func (s *CouponService) Validate(ctx context.Context, code string, orderAmount decimal.Decimal, now time.Time) (*dto.CouponResponse, error) {

	code = strings.TrimSpace(code)

	if code == "" {
		return nil, apperror.BadRequest(
			"COUPON_CODE_REQUIRED",
			"coupon code is required",
		)
	}

	if !orderAmount.GreaterThanOrEqual(decimal.Zero) {
		return nil, apperror.BadRequest(
			"INVALID_ORDER_AMOUNT",
			"order amount cannot be negative",
		)
	}

	coupon, err := s.couponRepo.GetByCode(ctx, code)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NotFound(
				"COUPON_NOT_FOUND",
				"coupon not found",
			)
		}

		return nil, err
	}

	if err := validateCoupon(coupon, orderAmount, now); err != nil {
		return nil, err
	}

	return toCouponResponse(coupon), nil
}

func validateCoupon(coupon *models.Coupon, orderAmount decimal.Decimal, now time.Time) error {

	if !coupon.IsActive {
		return apperror.BadRequest(
			"COUPON_INACTIVE",
			"coupon is inactive",
		)
	}

	if now.Before(coupon.StartsAt) {
		return apperror.BadRequest(
			"COUPON_NOT_STARTED",
			"coupon is not active yet",
		)
	}

	if !now.Before(coupon.ExpiresAt) {
		return apperror.BadRequest(
			"COUPON_EXPIRED",
			"coupon has expired",
		)
	}

	if coupon.UsageLimit != nil &&
		coupon.UsedCount >= *coupon.UsageLimit {
		return apperror.BadRequest(
			"COUPON_USAGE_LIMIT_REACHED",
			"coupon usage limit has been reached",
		)
	}

	if coupon.MinimumOrderAmount != nil &&
		orderAmount.LessThan(*coupon.MinimumOrderAmount) {
		return apperror.BadRequest(
			"MINIMUM_ORDER_AMOUNT_NOT_MET",
			"order amount does not meet the coupon minimum",
		)
	}

	return nil
}

func toCouponResponse(coupon *models.Coupon) *dto.CouponResponse {
	return &dto.CouponResponse{
		ID:                 coupon.ID,
		Code:               coupon.Code,
		DiscountType:       coupon.DiscountType,
		DiscountValue:      coupon.DiscountValue,
		MinimumOrderAmount: coupon.MinimumOrderAmount,
		MaximumDiscount:    coupon.MaximumDiscount,
		UsageLimit:         coupon.UsageLimit,
		UsedCount:          coupon.UsedCount,
		StartsAt:           coupon.StartsAt,
		ExpiresAt:          coupon.ExpiresAt,
		IsActive:           coupon.IsActive,
		CreatedAt:          coupon.CreatedAt,
		UpdatedAt:          coupon.UpdatedAt,
	}
}
