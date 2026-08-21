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

type ProductImageService struct {
	productImageRepo   *repositories.ProductImageRepository
	productVariantRepo *repositories.ProductVariantRepository
	productRepo        *repositories.ProductRepository
}

func NewProductImageService(productImageRepo *repositories.ProductImageRepository, productVariantRepo *repositories.ProductVariantRepository, productRepo *repositories.ProductRepository) *ProductImageService {
	return &ProductImageService{
		productImageRepo:   productImageRepo,
		productVariantRepo: productVariantRepo,
		productRepo:        productRepo,
	}
}

func (s *ProductImageService) Create(ctx context.Context, req *dto.CreateProductImageRequest) (*dto.ProductImageResponse, error) {
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

	if req.VariantID != nil {
		variant, err := s.productVariantRepo.GetByID(ctx, *req.VariantID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, apperror.NotFound(
					"PRODUCT_VARIANT_NOT_FOUND",
					"product variant not found",
				)
			}
			return nil, err
		}
		if variant.ProductID != req.ProductID {
			return nil, apperror.BadRequest(
				"INVALID_PRODUCT_VARIANT",
				"product variant does not belong to product",
			)
		}
	}

	if req.ImageURL == "" {
		return nil, apperror.BadRequest(
			"IMAGE_URL_REQUIRED",
			"image url required",
		)
	}

	if req.SortOrder < 0 {
		return nil, apperror.BadRequest(
			"INVALID_SORT_ORDER",
			"invalid sort order",
		)
	}

	image := &models.ProductImage{
		ProductID: req.ProductID,
		VariantID: req.VariantID,
		ImageURL:  req.ImageURL,
		IsPrimary: req.IsPrimary,
		SortOrder: req.SortOrder,
	}

	if err := s.productImageRepo.Create(ctx, image); err != nil {
		return nil, err
	}

	return toProductImageResponse(image), nil

}

func (s *ProductImageService) GetByID(ctx context.Context, id int64) (*dto.ProductImageResponse, error) {
	if id <= 0 {
		return nil, apperror.BadRequest(
			"INVALID_PRODUCT_IMAGE_ID",
			"invalid product image id",
		)
	}

	image, err := s.productImageRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NotFound(
				"PRODUCT_IMAGE_NOT_FOUND",
				"product image not found",
			)
		}
		return nil, err
	}
	return toProductImageResponse(image), nil

}

func (s *ProductImageService) GetAllByProductID(ctx context.Context, productID int64) ([]*dto.ProductImageResponse, error) {
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

	images, err := s.productImageRepo.GetAllByProductID(ctx, productID)
	if err != nil {
		return nil, err
	}

	result := make([]*dto.ProductImageResponse, 0, len(images))

	for _, image := range images {
		result = append(result, toProductImageResponse(image))
	}

	return result, nil

}

func (s *ProductImageService) Update(ctx context.Context, id int64, req *dto.UpdateProductImageRequest) (*dto.ProductImageResponse, error) {
	if id <= 0 {
		return nil, apperror.BadRequest(
			"INVALID_PRODUCT_IMAGE_ID",
			"invalid product image id",
		)
	}

	image, err := s.productImageRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NotFound(
				"PRODUCT_IMAGE_NOT_FOUND",
				"product image not found",
			)
		}

		return nil, err
	}

	if req.ImageURL != nil {
		if *req.ImageURL == "" {
			return nil, apperror.BadRequest(
				"IMAGE_URL_REQUIRED",
				"image url is required",
			)
		}

		image.ImageURL = *req.ImageURL
	}
	if req.VariantID != nil {
		variant, err := s.productVariantRepo.GetByID(ctx, *req.VariantID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, apperror.NotFound(
					"PRODUCT_VARIANT_NOT_FOUND",
					"product variant not found",
				)
			}
			return nil, err
		}

		if variant.ProductID != image.ProductID {
			return nil, apperror.BadRequest(
				"INVALID_PRODUCT_VARIANT",
				"product variant does not belong to product",
			)
		}

		image.VariantID = req.VariantID
	}

	if req.IsPrimary != nil {
		image.IsPrimary = *req.IsPrimary
	}

	if req.SortOrder != nil {
		if *req.SortOrder < 0 {
			return nil, apperror.BadRequest(
				"INVALID_SORT_ORDER",
				"sort order cannot be negative",
			)
		}

		image.SortOrder = *req.SortOrder
	}

	if err := s.productImageRepo.Update(ctx, image); err != nil {
		return nil, err
	}

	return toProductImageResponse(image), nil
}

func (s *ProductImageService) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return apperror.BadRequest(
			"INVALID_PRODUCT_IMAGE_ID",
			"invalid product image id",
		)
	}

	_, err := s.productImageRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return apperror.NotFound(
				"PRODUCT_IMAGE_NOT_FOUND",
				"product image not found",
			)
		}
		return err
	}
	return s.productImageRepo.Delete(ctx, id)

}

func toProductImageResponse(image *models.ProductImage) *dto.ProductImageResponse {
	return &dto.ProductImageResponse{
		ID:        image.ID,
		ProductID: image.ProductID,
		VariantID: image.VariantID,
		ImageURL:  image.ImageURL,
		IsPrimary: image.IsPrimary,
		SortOrder: image.SortOrder,
	}
}
