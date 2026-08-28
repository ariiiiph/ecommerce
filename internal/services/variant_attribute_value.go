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

type VariantAttributeValueService struct {
	variantAttributeValueRepo *repositories.VariantAttributeValueRepository
	productVariantRepo        *repositories.ProductVariantRepository
	attributeValueRepo        *repositories.AttributeValueRepository
}

func NewVariantAttributeValueService(variantAttributeValueRepo *repositories.VariantAttributeValueRepository, productVariantRepo *repositories.ProductVariantRepository, attributeValueRepo *repositories.AttributeValueRepository) *VariantAttributeValueService {
	return &VariantAttributeValueService{
		variantAttributeValueRepo: variantAttributeValueRepo,
		productVariantRepo:        productVariantRepo,
		attributeValueRepo:        attributeValueRepo,
	}
}

func (s *VariantAttributeValueService) Create(ctx context.Context, req *dto.CreateVariantAttributeValueRequest) (*dto.VariantAttributeValueResponse, error) {

	if req.VariantID <= 0 {
		return nil, apperror.BadRequest(
			"INVALID_VARIANT_ID",
			"invalid variant id",
		)
	}

	if req.AttributeValueID <= 0 {
		return nil, apperror.BadRequest(
			"INVALID_ATTRIBUTE_VALUE_ID",
			"invalid attribute value id",
		)
	}

	_, err := s.productVariantRepo.GetByID(ctx, req.VariantID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NotFound(
				"PRODUCT_VARIANT_NOT_FOUND",
				"product variant not found",
			)
		}

		return nil, err
	}

	_, err = s.attributeValueRepo.GetByID(ctx, req.AttributeValueID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NotFound(
				"ATTRIBUTE_VALUE_NOT_FOUND",
				"attribute value not found",
			)
		}

		return nil, err
	}

	relationship := &models.VariantAttributeValue{
		VariantID:        req.VariantID,
		AttributeValueID: req.AttributeValueID,
	}

	if err := s.variantAttributeValueRepo.Create(ctx, relationship); err != nil {
		return nil, err
	}

	return toVariantAttributeValueResponse(relationship), nil
}

func (s *VariantAttributeValueService) GetByVariantID(ctx context.Context, variantID int64) ([]*dto.VariantAttributeValueResponse, error) {

	if variantID <= 0 {
		return nil, apperror.BadRequest(
			"INVALID_VARIANT_ID",
			"invalid variant id",
		)
	}

	_, err := s.productVariantRepo.GetByID(ctx, variantID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NotFound(
				"PRODUCT_VARIANT_NOT_FOUND",
				"product variant not found",
			)
		}

		return nil, err
	}

	relationships, err := s.variantAttributeValueRepo.GetByVariantID(
		ctx,
		variantID,
	)
	if err != nil {
		return nil, err
	}

	result := make(
		[]*dto.VariantAttributeValueResponse,
		0,
		len(relationships),
	)

	for _, relationship := range relationships {
		result = append(
			result,
			toVariantAttributeValueResponse(relationship),
		)
	}

	return result, nil
}

func (s *VariantAttributeValueService) Delete(ctx context.Context, variantID int64, attributeValueID int64) error {

	if variantID <= 0 {
		return apperror.BadRequest(
			"INVALID_VARIANT_ID",
			"invalid variant id",
		)
	}

	if attributeValueID <= 0 {
		return apperror.BadRequest(
			"INVALID_ATTRIBUTE_VALUE_ID",
			"invalid attribute value id",
		)
	}

	_, err := s.productVariantRepo.GetByID(ctx, variantID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return apperror.NotFound(
				"PRODUCT_VARIANT_NOT_FOUND",
				"product variant not found",
			)
		}

		return err
	}

	_, err = s.attributeValueRepo.GetByID(ctx, attributeValueID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return apperror.NotFound(
				"ATTRIBUTE_VALUE_NOT_FOUND",
				"attribute value not found",
			)
		}

		return err
	}

	if err := s.variantAttributeValueRepo.Delete(
		ctx,
		variantID,
		attributeValueID,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return apperror.NotFound(
				"VARIANT_ATTRIBUTE_VALUE_NOT_FOUND",
				"variant attribute value not found",
			)
		}

		return err
	}

	return nil
}

func toVariantAttributeValueResponse(relationship *models.VariantAttributeValue) *dto.VariantAttributeValueResponse {
	return &dto.VariantAttributeValueResponse{
		VariantID:        relationship.VariantID,
		AttributeValueID: relationship.AttributeValueID,
	}
}
