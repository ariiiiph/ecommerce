package services

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/ariiiiph/ecommerce/internal/dto"
	"github.com/ariiiiph/ecommerce/internal/models"
	"github.com/ariiiiph/ecommerce/internal/repositories"
)

var (
	ErrBrandNotFound      = errors.New("brand not found")
	ErrBrandAlreadyExists = errors.New("brand already exists")
)

type BrandService struct {
	brandRepo *repositories.BrandRepository
}

func NewBrandService(brandRepo *repositories.BrandRepository) *BrandService {
	return &BrandService{
		brandRepo: brandRepo,
	}
}

func (s *BrandService) Create(ctx context.Context, req *dto.CreateBrandRequest) (*dto.BrandResponse, error) {

	if req.Name == "" {
		return nil, errors.New("brand name is required")
	}

	if req.Slug == "" {
		return nil, errors.New("brand slug is required")
	}

	brand := &models.Brand{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
	}

	if err := s.brandRepo.Create(ctx, brand); err != nil {
		return nil, err
	}

	return &dto.BrandResponse{
		ID:          brand.ID,
		Name:        brand.Name,
		Slug:        brand.Slug,
		Description: brand.Description,
	}, nil
}

func (s *BrandService) GetByID(ctx context.Context, id int64) (*dto.BrandResponse, error) {
	brand, err := s.brandRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrBrandNotFound
		}
		return nil, err
	}

	return &dto.BrandResponse{
		ID:          brand.ID,
		Name:        brand.Name,
		Slug:        brand.Slug,
		Description: brand.Description,
	}, nil
}

func (s *BrandService) GetAll(ctx context.Context) ([]*dto.BrandResponse, error) {
	brands, err := s.brandRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]*dto.BrandResponse, 0, len(brands))

	for _, brand := range brands {
		result = append(result, &dto.BrandResponse{
			ID:          brand.ID,
			Name:        brand.Name,
			Slug:        brand.Slug,
			Description: brand.Description,
		})
	}
	return result, nil

}

func (s *BrandService) Update(ctx context.Context, id int64, req *dto.UpdateBrandRequest) (*dto.BrandResponse, error) {
	brand, err := s.brandRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrBrandNotFound
		}
		return nil, err
	}

	if req.Name == "" {
		return nil, errors.New("brand name is required")
	}

	if req.Slug == "" {
		return nil, errors.New("brand slug is required")
	}

	brand.Name = req.Name
	brand.Slug = req.Slug
	brand.Description = req.Description
	brand.UpdatedAt = time.Now()

	if err := s.brandRepo.Update(ctx, brand); err != nil {
		return nil, err
	}

	return &dto.BrandResponse{
		ID:          brand.ID,
		Name:        brand.Name,
		Slug:        brand.Slug,
		Description: brand.Description,
	}, nil
}

func (s *BrandService) Delete(ctx context.Context, id int64) error {
	_, err := s.brandRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrBrandNotFound
		}

		return err
	}

	return s.brandRepo.Delete(ctx, id)
}
