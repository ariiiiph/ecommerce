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

type InventoryService struct {
	inventoryRepo *repositories.InventoryRepository
	variantRepo   *repositories.ProductVariantRepository
}

func NewInventoryService(inventoryRepo *repositories.InventoryRepository, variantRepo *repositories.ProductVariantRepository) *InventoryService {
	return &InventoryService{
		inventoryRepo: inventoryRepo,
		variantRepo:   variantRepo,
	}
}

func (s *InventoryService) Create(ctx context.Context, req *dto.CreateInventoryRequest) (*dto.InventoryResponse, error) {
	if req.VariantID <= 0 {
		return nil, apperror.BadRequest(
			"INVALID_VARIANT_ID",
			"invalid variant id",
		)
	}

	if req.Quantity < 0 {
		return nil, apperror.BadRequest(
			"INVALID_QUANTITY",
			"quantity can not be negative",
		)
	}

	if req.ReservedQuantity < 0 {
		return nil, apperror.BadRequest(
			"INVALID_RESERVED_QUANTITY",
			"reserved quantity can not be negative",
		)
	}
	if req.LowStockThreshold < 0 {
		return nil, apperror.BadRequest(
			"INVALID_LOW_STOCK_THRESHOLD",
			"low stock threshold cannot be negative",
		)
	}

	if req.ReservedQuantity > req.Quantity {
		return nil, apperror.BadRequest(
			"INVALID_RESERVED_QUANTITY",
			"reserved quantity cannot be greater than quantity",
		)
	}

	_, err := s.variantRepo.GetByID(ctx, req.VariantID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NotFound(
				"VARIANT_NOT_FOUND",
				"product variant not found",
			)
		}
		return nil, err
	}
	inventory := &models.Inventory{
		VariantID:         req.VariantID,
		Quantity:          req.Quantity,
		ReservedQuantity:  req.ReservedQuantity,
		LowStockThreshold: req.LowStockThreshold,
	}

	if err := s.inventoryRepo.Create(ctx, inventory); err != nil {
		return nil, err
	}

	return toInventoryResponse(inventory), nil

}

func (s *InventoryService) GetByVariantID(ctx context.Context, variantID int64) (*dto.InventoryResponse, error) {
	if variantID <= 0 {
		return nil, apperror.BadRequest(
			"INVALID_VARIANT_ID",
			"invalid variant id",
		)
	}

	inventory, err := s.inventoryRepo.GetByVariantID(ctx, variantID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NotFound(
				"INVENTORY_NOT_FOUND",
				"inventory not found",
			)
		}
		return nil, err
	}

	return toInventoryResponse(inventory), nil

}

func (s *InventoryService) GetAll(ctx context.Context) ([]*dto.InventoryResponse, error) {
	inventories, err := s.inventoryRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]*dto.InventoryResponse, 0, len(inventories))

	for _, inventory := range inventories {
		result = append(result, toInventoryResponse(inventory))
	}
	return result, nil
}

func (s *InventoryService) Update(ctx context.Context, variantID int64, req *dto.UpdateInventoryRequest) (*dto.InventoryResponse, error) {
	if variantID <= 0 {
		return nil, apperror.BadRequest(
			"INVALID_VARIANT_ID",
			"invalid variant id",
		)
	}
	inventory, err := s.inventoryRepo.GetByVariantID(ctx, variantID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NotFound(
				"INVENTORY_NOT_FOUND",
				"inventory not found",
			)
		}
		return nil, err
	}

	if req.Quantity != nil {
		if *req.Quantity < 0 {
			return nil, apperror.BadRequest(
				"INVALID_QUANTITY",
				"quantity can not be negative",
			)
		}
		inventory.Quantity = *req.Quantity
	}
	if req.ReservedQuantity != nil {
		if *req.ReservedQuantity < 0 {
			return nil, apperror.BadRequest(
				"INVALID_RESERVED_QUANTITY",
				"reserved quantity cannot be negative",
			)
		}

		inventory.ReservedQuantity = *req.ReservedQuantity
	}

	if req.LowStockThreshold != nil {
		if *req.LowStockThreshold < 0 {
			return nil, apperror.BadRequest(
				"INVALID_LOW_STOCK_THRESHOLD",
				"low stock threshold cannot be negative",
			)
		}

		inventory.LowStockThreshold = *req.LowStockThreshold
	}

	if inventory.ReservedQuantity > inventory.Quantity {
		return nil, apperror.BadRequest(
			"INVALID_RESERVED_QUANTITY",
			"reserved quantity cannot be greater than quantity",
		)
	}

	if err := s.inventoryRepo.Update(ctx, inventory); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NotFound(
				"INVENTORY_NOT_FOUND",
				"inventory not found",
			)
		}

		return nil, err
	}

	return toInventoryResponse(inventory), nil

}

func (s *InventoryService) Delete(ctx context.Context, variantID int64) error {
	if variantID <= 0 {
		return apperror.BadRequest(
			"INVALID_VARIANT_ID",
			"invalid variant id",
		)
	}

	if _, err := s.inventoryRepo.GetByVariantID(ctx, variantID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return apperror.NotFound(
				"INVENTORY_NOT_FOUND",
				"inventory not found",
			)
		}

		return err
	}

	return s.inventoryRepo.Delete(ctx, variantID)

}

func toInventoryResponse(inventory *models.Inventory) *dto.InventoryResponse {
	return &dto.InventoryResponse{
		VariantID:         inventory.VariantID,
		Quantity:          inventory.Quantity,
		ReservedQuantity:  inventory.ReservedQuantity,
		LowStockThreshold: inventory.LowStockThreshold,
		CreatedAt:         inventory.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:         inventory.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
