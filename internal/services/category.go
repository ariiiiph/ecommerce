package services

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/ariiiiph/ecommerce/internal/apperror"
	"github.com/ariiiiph/ecommerce/internal/dto"
	"github.com/ariiiiph/ecommerce/internal/models"
	"github.com/ariiiiph/ecommerce/internal/repositories"
)

type CategoryService struct {
	categoryRepo *repositories.CategoryRepository
}

func NewCategoryService(categoryRepo *repositories.CategoryRepository) *CategoryService {
	return &CategoryService{
		categoryRepo: categoryRepo,
	}
}

func (s *CategoryService) Create(ctx context.Context, req *dto.CreateCategoryRequest) (*dto.CategoryResponse, error) {
	if req.Name == "" {
		return nil, apperror.BadRequest(
			"CATEGORY_NAME_REQUIRED",
			"category name is required",
		)
	}

	if req.Slug == "" {
		return nil, apperror.BadRequest(
			"CATEGORY_SLUG_REQUIRED",
			"category slug is required",
		)
	}

	if req.ParentID != nil {
		_, err := s.categoryRepo.FindByID(ctx, *req.ParentID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, apperror.NotFound(
					"PARENT_CATEGORY_NOT_FOUND",
					"parent category not found",
				)
			}

			return nil, err
		}
	}

	category := &models.Category{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		ParentID:    req.ParentID,
	}

	if err := s.categoryRepo.Create(ctx, category); err != nil {
		return nil, err
	}

	return &dto.CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Slug:        category.Slug,
		Description: category.Description,
		ParentID:    category.ParentID,
	}, nil
}

func (s *CategoryService) GetByID(ctx context.Context, id int64) (*dto.CategoryResponse, error) {
	category, err := s.categoryRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NotFound(
				"CATEGORY_NOT_FOUND",
				"category not found",
			)
		}
		return nil, err
	}
	return &dto.CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Slug:        category.Slug,
		Description: category.Description,
		ParentID:    category.ParentID,
	}, nil
}

func (s *CategoryService) GetAll(ctx context.Context) ([]*dto.CategoryResponse, error) {
	categories, err := s.categoryRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]*dto.CategoryResponse, 0, len(categories))

	for _, category := range categories {
		result = append(result, &dto.CategoryResponse{
			ID:          category.ID,
			Name:        category.Name,
			Slug:        category.Slug,
			Description: category.Description,
			ParentID:    category.ParentID,
		})
	}

	return result, nil
}

func (s *CategoryService) Update(ctx context.Context, id int64, req *dto.UpdateCategoryRequest) (*dto.CategoryResponse, error) {
	category, err := s.categoryRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NotFound(
				"CATEGORY_NOT_FOUND",
				"category not found",
			)
		}

		return nil, err
	}

	if req.Name == "" {
		return nil, apperror.BadRequest(
			"CATEGORY_NAME_REQUIRED",
			"category name is required",
		)
	}

	if req.Slug == "" {
		return nil, apperror.BadRequest(
			"CATEGORY_SLUG_REQUIRED",
			"category slug is required",
		)
	}

	if req.ParentID != nil {
		if *req.ParentID == id {
			return nil, apperror.BadRequest(
				"CATEGORY_SELF_PARENT",
				"category cannot be its own parent",
			)
		}

		_, err := s.categoryRepo.FindByID(ctx, *req.ParentID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, apperror.NotFound(
					"PARENT_CATEGORY_NOT_FOUND",
					"parent category not found",
				)
			}

			return nil, err
		}
	}

	category.Name = req.Name
	category.Slug = req.Slug
	category.Description = req.Description
	category.ParentID = req.ParentID
	category.UpdatedAt = time.Now()

	if err := s.categoryRepo.Update(ctx, category); err != nil {
		return nil, err
	}

	return &dto.CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Slug:        category.Slug,
		Description: category.Description,
		ParentID:    category.ParentID,
	}, nil
}

func (s *CategoryService) Delete(ctx context.Context, id int64) error {
	_, err := s.categoryRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return apperror.NotFound(
				"CATEGORY_NOT_FOUND",
				"category not found",
			)
		}

		return err
	}

	return s.categoryRepo.Delete(ctx, id)
}

func (s *CategoryService) GetChildren(ctx context.Context, parentID int64) ([]*dto.CategoryResponse, error) {
	_, err := s.categoryRepo.FindByID(ctx, parentID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NotFound(
				"CATEGORY_NOT_FOUND",
				"category not found",
			)
		}

		return nil, err
	}
	categories, err := s.categoryRepo.FindChildren(ctx, parentID)
	if err != nil {
		return nil, err
	}

	result := make([]*dto.CategoryResponse, 0, len(categories))

	for _, category := range categories {
		result = append(result, &dto.CategoryResponse{
			ID:          category.ID,
			Name:        category.Name,
			Slug:        category.Slug,
			Description: category.Description,
			ParentID:    category.ParentID,
		})
	}

	return result, nil
}

func (s *CategoryService) GetTree(ctx context.Context, parentID int64) ([]*dto.CategoryResponse, error) {
	_, err := s.categoryRepo.FindByID(ctx, parentID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NotFound(
				"CATEGORY_NOT_FOUND",
				"category not found",
			)
		}

		return nil, err
	}

	categories, err := s.categoryRepo.FindTree(ctx, parentID)
	if err != nil {
		return nil, err
	}

	result := make([]*dto.CategoryResponse, 0, len(categories))

	for _, category := range categories {
		result = append(result, &dto.CategoryResponse{
			ID:          category.ID,
			Name:        category.Name,
			Slug:        category.Slug,
			Description: category.Description,
			ParentID:    category.ParentID,
		})
	}

	return result, nil
}
