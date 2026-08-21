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

type AttributeService struct {
	attributeRepo *repositories.AttributeRepository
}

func NewAttributeService(attributeRepo *repositories.AttributeRepository) *AttributeService {
	return &AttributeService{
		attributeRepo: attributeRepo,
	}
}

func (s *AttributeService) Create(ctx context.Context, req *dto.CreateAttributeRequest) (*dto.AttributeResponse, error) {

	name := strings.TrimSpace(req.Name)

	if name == "" {
		return nil, apperror.BadRequest(
			"ATTRIBUTE_NAME_REQUIRED",
			"attribute name is required",
		)
	}

	attribute := &models.Attribute{
		Name: name,
	}

	if err := s.attributeRepo.Create(ctx, attribute); err != nil {
		return nil, err
	}

	return toAttributeResponse(attribute), nil
}

func (s *AttributeService) GetByID(ctx context.Context, id int64) (*dto.AttributeResponse, error) {
	if id <= 0 {
		return nil, apperror.BadRequest(
			"INVALID_ATTRIBUTE_ID",
			"invalid attribute id",
		)
	}

	attribute, err := s.attributeRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NotFound(
				"ATTRIBUTE_NOT_FOUND",
				"attribute not found",
			)
		}

		return nil, err
	}

	return toAttributeResponse(attribute), nil
}

func (s *AttributeService) GetAll(ctx context.Context) ([]*dto.AttributeResponse, error) {
	attributes, err := s.attributeRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]*dto.AttributeResponse, 0, len(attributes))

	for _, attribute := range attributes {
		result = append(result, toAttributeResponse(attribute))
	}

	return result, nil
}

func (s *AttributeService) Update(ctx context.Context, id int64, req *dto.UpdateAttributeRequest) (*dto.AttributeResponse, error) {
	if id <= 0 {
		return nil, apperror.BadRequest(
			"INVALID_ATTRIBUTE_ID",
			"invalid attribute id",
		)
	}

	attribute, err := s.attributeRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NotFound(
				"ATTRIBUTE_NOT_FOUND",
				"attribute not found",
			)
		}

		return nil, err
	}

	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)

		if name == "" {
			return nil, apperror.BadRequest(
				"ATTRIBUTE_NAME_REQUIRED",
				"attribute name is required",
			)
		}

		attribute.Name = name
	}

	if req.IsActive != nil {
		attribute.IsActive = *req.IsActive
	}

	if err := s.attributeRepo.Update(ctx, attribute); err != nil {
		return nil, err
	}

	return toAttributeResponse(attribute), nil
}

func (s *AttributeService) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return apperror.BadRequest(
			"INVALID_ATTRIBUTE_ID",
			"invalid attribute id",
		)
	}

	_, err := s.attributeRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return apperror.NotFound(
				"ATTRIBUTE_NOT_FOUND",
				"attribute not found",
			)
		}

		return err
	}

	return s.attributeRepo.Delete(ctx, id)
}

func toAttributeResponse(attribute *models.Attribute) *dto.AttributeResponse {
	return &dto.AttributeResponse{
		ID:        attribute.ID,
		Name:      attribute.Name,
		IsActive:  attribute.IsActive,
		CreatedAt: attribute.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: attribute.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
