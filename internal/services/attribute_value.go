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

type AttributeValueService struct {
	attributeValueRepo *repositories.AttributeValueRepository
	attributeRepo      *repositories.AttributeRepository
}

func NewAttributeValueService(attributeValueRepo *repositories.AttributeValueRepository, attributeRepo *repositories.AttributeRepository) *AttributeValueService {
	return &AttributeValueService{
		attributeValueRepo: attributeValueRepo,
		attributeRepo:      attributeRepo,
	}
}

func (s *AttributeValueService) Create(ctx context.Context, req *dto.CreateAttributeValueRequest) (*dto.AttributeValueResponse, error) {

	if req.AttributeID <= 0 {
		return nil, apperror.BadRequest(
			"INVALID_ATTRIBUTE_ID",
			"invalid attribute id",
		)
	}

	value := strings.TrimSpace(req.Value)

	if value == "" {
		return nil, apperror.BadRequest(
			"ATTRIBUTE_VALUE_REQUIRED",
			"attribute value is required",
		)
	}

	_, err := s.attributeRepo.GetByID(ctx, req.AttributeID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NotFound(
				"ATTRIBUTE_NOT_FOUND",
				"attribute not found",
			)
		}

		return nil, err
	}

	attributeValue := &models.AttributeValue{
		AttributeID: req.AttributeID,
		Value:       value,
	}

	if err := s.attributeValueRepo.Create(ctx, attributeValue); err != nil {
		return nil, err
	}

	return toAttributeValueResponse(attributeValue), nil
}

func (s *AttributeValueService) GetByID(ctx context.Context, id int64) (*dto.AttributeValueResponse, error) {

	if id <= 0 {
		return nil, apperror.BadRequest(
			"INVALID_ATTRIBUTE_VALUE_ID",
			"invalid attribute value id",
		)
	}

	attributeValue, err := s.attributeValueRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NotFound(
				"ATTRIBUTE_VALUE_NOT_FOUND",
				"attribute value not found",
			)
		}

		return nil, err
	}

	return toAttributeValueResponse(attributeValue), nil
}

func (s *AttributeValueService) GetAll(ctx context.Context) ([]*dto.AttributeValueResponse, error) {

	attributeValues, err := s.attributeValueRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]*dto.AttributeValueResponse, 0, len(attributeValues))

	for _, attributeValue := range attributeValues {
		result = append(
			result,
			toAttributeValueResponse(attributeValue),
		)
	}

	return result, nil
}

func (s *AttributeValueService) Update(ctx context.Context, id int64, req *dto.UpdateAttributeValueRequest) (*dto.AttributeValueResponse, error) {

	if id <= 0 {
		return nil, apperror.BadRequest(
			"INVALID_ATTRIBUTE_VALUE_ID",
			"invalid attribute value id",
		)
	}

	attributeValue, err := s.attributeValueRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NotFound(
				"ATTRIBUTE_VALUE_NOT_FOUND",
				"attribute value not found",
			)
		}

		return nil, err
	}

	if req.Value != nil {
		value := strings.TrimSpace(*req.Value)

		if value == "" {
			return nil, apperror.BadRequest(
				"ATTRIBUTE_VALUE_REQUIRED",
				"attribute value is required",
			)
		}

		attributeValue.Value = value
	}

	if err := s.attributeValueRepo.Update(ctx, attributeValue); err != nil {
		return nil, err
	}

	return toAttributeValueResponse(attributeValue), nil
}

func (s *AttributeValueService) Delete(ctx context.Context, id int64) error {

	if id <= 0 {
		return apperror.BadRequest(
			"INVALID_ATTRIBUTE_VALUE_ID",
			"invalid attribute value id",
		)
	}

	_, err := s.attributeValueRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return apperror.NotFound(
				"ATTRIBUTE_VALUE_NOT_FOUND",
				"attribute value not found",
			)
		}

		return err
	}

	return s.attributeValueRepo.Delete(ctx, id)
}

func toAttributeValueResponse(attributeValue *models.AttributeValue) *dto.AttributeValueResponse {
	return &dto.AttributeValueResponse{
		ID:          attributeValue.ID,
		AttributeID: attributeValue.AttributeID,
		Value:       attributeValue.Value,
		CreatedAt:   attributeValue.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
