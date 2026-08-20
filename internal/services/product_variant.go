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

type ProductVariantService struct {
	productVariantRepo *repositories.ProductVariantRepository
	productRepo        *repositories.ProductRepository
}

func NewProductVariantService(productVariantRepo *repositories.ProductVariantRepository, productRepo *repositories.ProductRepository) *ProductVariantService {
	return &ProductVariantService{
		productVariantRepo: productVariantRepo,
		productRepo:        productRepo,
	}
}

func (s *ProductVariantService) Create(ctx context.Context, req *dto.CreateProductVariantRequest) (*dto.ProductVariantResponse, error) {

	if req.ProductID <= 0 {
		return nil, apperror.BadRequest(
			"INVALID_PRODUCT_ID",
			"invalid product id",
		)
	}

	if req.SKU == "" {
		return nil, apperror.BadRequest(
			"PRODUCT_VARIANT_SKU_REQUIRED",
			"sku is required",
		)
	}

	if req.Price < 0 {
		return nil, apperror.BadRequest(
			"INVALID_PRODUCT_VARIANT_PRICE",
			"price cannot be negative",
		)
	}

	if req.DiscountPrice != nil {
		if *req.DiscountPrice < 0 {
			return nil, apperror.BadRequest(
				"INVALID_PRODUCT_VARIANT_DISCOUNT_PRICE",
				"discount price cannot be negative",
			)
		}

		if *req.DiscountPrice > req.Price {
			return nil, apperror.BadRequest(
				"INVALID_PRODUCT_VARIANT_DISCOUNT_PRICE",
				"discount price cannot be greater than price",
			)
		}
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

	variant := &models.ProductVariant{
		ProductID:     req.ProductID,
		SKU:           req.SKU,
		Price:         req.Price,
		DiscountPrice: req.DiscountPrice,
	}

	if err := s.productVariantRepo.Create(ctx, variant); err != nil {
		return nil, err
	}

	return toProductVariantResponse(variant), nil
}

func (s *ProductVariantService) GetByID(ctx context.Context, id int64) (*dto.ProductVariantResponse, error) {

	if id <= 0 {
		return nil, apperror.BadRequest(
			"INVALID_PRODUCT_VARIANT_ID",
			"invalid product variant id",
		)
	}

	variant, err := s.productVariantRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NotFound(
				"PRODUCT_VARIANT_NOT_FOUND",
				"product variant not found",
			)
		}

		return nil, err
	}

	return toProductVariantResponse(variant), nil
}

func (s *ProductVariantService) GetAllByProductID(ctx context.Context, productID int64) ([]*dto.ProductVariantResponse, error) {

	if productID <= 0 {
		return nil, apperror.BadRequest(
			"INVALID_PRODUCT_ID",
			"invalid product id",
		)
	}

	_, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NotFound(
				"PRODUCT_NOT_FOUND",
				"product not found",
			)
		}

		return nil, err
	}

	variants, err := s.productVariantRepo.GetAllByProductID(ctx, productID)
	if err != nil {
		return nil, err
	}

	result := make([]*dto.ProductVariantResponse, 0, len(variants))

	for _, variant := range variants {
		result = append(result, toProductVariantResponse(variant))
	}

	return result, nil
}

func (s *ProductVariantService) Update(ctx context.Context, id int64, req *dto.UpdateProductVariantRequest) (*dto.ProductVariantResponse, error) {

	if id <= 0 {
		return nil, apperror.BadRequest(
			"INVALID_PRODUCT_VARIANT_ID",
			"invalid product variant id",
		)
	}

	variant, err := s.productVariantRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NotFound(
				"PRODUCT_VARIANT_NOT_FOUND",
				"product variant not found",
			)
		}

		return nil, err
	}

	if req.SKU != nil {
		if *req.SKU == "" {
			return nil, apperror.BadRequest(
				"PRODUCT_VARIANT_SKU_REQUIRED",
				"sku is required",
			)
		}

		variant.SKU = *req.SKU
	}

	if req.Price != nil {
		if *req.Price < 0 {
			return nil, apperror.BadRequest(
				"INVALID_PRODUCT_VARIANT_PRICE",
				"price cannot be negative",
			)
		}

		variant.Price = *req.Price
	}

	if req.DiscountPrice != nil {
		if *req.DiscountPrice < 0 {
			return nil, apperror.BadRequest(
				"INVALID_PRODUCT_VARIANT_DISCOUNT_PRICE",
				"discount price cannot be negative",
			)
		}

		if *req.DiscountPrice > variant.Price {
			return nil, apperror.BadRequest(
				"INVALID_PRODUCT_VARIANT_DISCOUNT_PRICE",
				"discount price cannot be greater than price",
			)
		}

		variant.DiscountPrice = req.DiscountPrice
	}

	if variant.DiscountPrice != nil && *variant.DiscountPrice > variant.Price {
		return nil, apperror.BadRequest(
			"INVALID_PRODUCT_VARIANT_DISCOUNT_PRICE",
			"discount price cannot be greater than price",
		)
	}

	if err := s.productVariantRepo.Update(ctx, variant); err != nil {
		return nil, err
	}

	return toProductVariantResponse(variant), nil
}

func (s *ProductVariantService) Delete(ctx context.Context, id int64) error {

	if id <= 0 {
		return apperror.BadRequest(
			"INVALID_PRODUCT_VARIANT_ID",
			"invalid product variant id",
		)
	}

	_, err := s.productVariantRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return apperror.NotFound(
				"PRODUCT_VARIANT_NOT_FOUND",
				"product variant not found",
			)
		}

		return err
	}

	return s.productVariantRepo.Delete(ctx, id)
}

func toProductVariantResponse(variant *models.ProductVariant) *dto.ProductVariantResponse {
	return &dto.ProductVariantResponse{
		ID:            variant.ID,
		ProductID:     variant.ProductID,
		SKU:           variant.SKU,
		Price:         variant.Price,
		DiscountPrice: variant.DiscountPrice,
	}
}
