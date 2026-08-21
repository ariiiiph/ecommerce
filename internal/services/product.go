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

type ProductService struct {
	productRepo  *repositories.ProductRepository
	brandRepo    *repositories.BrandRepository
	categoryRepo *repositories.CategoryRepository
}

func NewProductService(
	productRepo *repositories.ProductRepository,
	brandRepo *repositories.BrandRepository,
	categoryRepo *repositories.CategoryRepository,
) *ProductService {
	return &ProductService{
		productRepo:  productRepo,
		brandRepo:    brandRepo,
		categoryRepo: categoryRepo,
	}
}

func (s *ProductService) Create(ctx context.Context, req *dto.CreateProductRequest) (*dto.ProductResponse, error) {

	if req.Name == "" {
		return nil, apperror.BadRequest(
			"PRODUCT_NAME_REQUIRED",
			"product name is required",
		)
	}

	if req.Slug == "" {
		return nil, apperror.BadRequest(
			"PRODUCT_SLUG_REQUIRED",
			"product slug is required",
		)
	}

	if !isValidProductStatus(req.Status) {
		return nil, apperror.BadRequest(
			"INVALID_PRODUCT_STATUS",
			"invalid product status",
		)
	}
	if req.BrandID <= 0 {
		return nil, apperror.BadRequest(
			"INVALID_BRAND_ID",
			"invalid brand id",
		)
	}

	if req.CategoryID <= 0 {
		return nil, apperror.BadRequest(
			"INVALID_CATEGORY_ID",
			"invalid category id",
		)
	}

	_, err := s.brandRepo.FindByID(ctx, req.BrandID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NotFound(
				"BRAND_NOT_FOUND",
				"brand not found",
			)
		}

		return nil, err
	}

	_, err = s.categoryRepo.FindByID(ctx, req.CategoryID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NotFound(
				"CATEGORY_NOT_FOUND",
				"category not found",
			)
		}

		return nil, err
	}

	product := &models.Product{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		BrandID:     req.BrandID,
		CategoryID:  req.CategoryID,
		Status:      req.Status,
	}

	if err := s.productRepo.Create(ctx, product); err != nil {
		return nil, err
	}

	return toProductResponse(product), nil
}

func (s *ProductService) GetByID(ctx context.Context, id int64) (*dto.ProductResponse, error) {

	product, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NotFound(
				"PRODUCT_NOT_FOUND",
				"product not found",
			)
		}

		return nil, err
	}

	return toProductResponse(product), nil
}

func (s *ProductService) GetAll(ctx context.Context, req *dto.PaginationRequest) ([]*dto.ProductResponse, *dto.PaginationResponse, error) {

	if req.Page <= 0 {
		req.Page = 1
	}

	if req.Limit <= 0 {
		req.Limit = 10
	}

	if req.Limit > 100 {
		req.Limit = 100
	}

	products, total, err := s.productRepo.GetAll(ctx, req.Page, req.Limit)
	if err != nil {
		return nil, nil, err
	}

	result := make([]*dto.ProductResponse, 0, len(products))

	for _, product := range products {
		result = append(result, toProductResponse(product))
	}
	totalPages := 0
	if total > 0 {
		totalPages = (total + req.Limit - 1) / req.Limit
	}
	pagination := &dto.PaginationResponse{
		Page:       req.Page,
		Limit:      req.Limit,
		Total:      total,
		TotalPages: totalPages,
	}

	return result, pagination, nil
}

func (s *ProductService) Update(ctx context.Context, id int64, req *dto.UpdateProductRequest) (*dto.ProductResponse, error) {

	product, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NotFound(
				"PRODUCT_NOT_FOUND",
				"product not found",
			)
		}

		return nil, err
	}

	if req.Name == "" {
		return nil, apperror.BadRequest(
			"PRODUCT_NAME_REQUIRED",
			"product name is required",
		)
	}

	if req.Slug == "" {
		return nil, apperror.BadRequest(
			"PRODUCT_SLUG_REQUIRED",
			"product slug is required",
		)
	}

	if !isValidProductStatus(req.Status) {
		return nil, apperror.BadRequest(
			"INVALID_PRODUCT_STATUS",
			"invalid product status",
		)
	}
	if req.BrandID <= 0 {
		return nil, apperror.BadRequest(
			"INVALID_BRAND_ID",
			"invalid brand id",
		)
	}

	if req.CategoryID <= 0 {
		return nil, apperror.BadRequest(
			"INVALID_CATEGORY_ID",
			"invalid category id",
		)
	}

	_, err = s.brandRepo.FindByID(ctx, req.BrandID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NotFound(
				"BRAND_NOT_FOUND",
				"brand not found",
			)
		}

		return nil, err
	}

	_, err = s.categoryRepo.FindByID(ctx, req.CategoryID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NotFound(
				"CATEGORY_NOT_FOUND",
				"category not found",
			)
		}

		return nil, err
	}

	product.Name = req.Name
	product.Slug = req.Slug
	product.Description = req.Description
	product.BrandID = req.BrandID
	product.CategoryID = req.CategoryID
	product.Status = req.Status

	if err := s.productRepo.Update(ctx, product); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NotFound(
				"PRODUCT_NOT_FOUND",
				"product not found",
			)
		}

		return nil, err
	}

	return toProductResponse(product), nil
}

func (s *ProductService) Delete(ctx context.Context, id int64) error {
	if _, err := s.productRepo.GetByID(ctx, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return apperror.NotFound(
				"PRODUCT_NOT_FOUND",
				"product not found",
			)
		}

		return err
	}

	return s.productRepo.Delete(ctx, id)
}

func isValidProductStatus(status string) bool {
	switch status {
	case "draft", "active", "inactive":
		return true
	default:
		return false
	}
}

func toProductResponse(product *models.Product) *dto.ProductResponse {
	return &dto.ProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		Slug:        product.Slug,
		Description: product.Description,
		BrandID:     product.BrandID,
		CategoryID:  product.CategoryID,
		Status:      product.Status,
	}
}
