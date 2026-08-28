package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ariiiiph/ecommerce/internal/apperror"
	"github.com/ariiiiph/ecommerce/internal/dto"
	"github.com/ariiiiph/ecommerce/internal/services"
)

type InventoryHandler struct {
	inventoryService *services.InventoryService
}

func NewInventoryHandler(inventoryService *services.InventoryService) *InventoryHandler {
	return &InventoryHandler{
		inventoryService: inventoryService,
	}
}

// Create godoc
// @Summary Create inventory
// @Description Creates inventory for a product variant.
// @Tags Inventory
// @Accept json
// @Produce json
// @Param request body dto.CreateInventoryRequest true "Inventory data"
// @Security BearerAuth
// @Success 201 {object} dto.InventoryResponse
// @Failure 400 {object} map[string]any
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Failure 404 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /api/inventory [post]
func (h *InventoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateInventoryRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.Write(
			w,
			apperror.BadRequest(
				"INVALID_REQUEST_BODY",
				"invalid request body",
			),
		)
		return
	}
	result, err := h.inventoryService.Create(r.Context(), &req)
	if err != nil {
		apperror.Write(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, result)

}

// GetByVariantID godoc
// @Summary Get inventory by variant ID
// @Description Retrieves inventory for a product variant.
// @Tags Inventory
// @Produce json
// @Param variant_id path int true "Product Variant ID"
// @Success 200 {object} dto.InventoryResponse
// @Failure 400 {object} map[string]any
// @Failure 404 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /api/inventory/{variant_id} [get]
func (h *InventoryHandler) GetByVariantID(w http.ResponseWriter, r *http.Request) {
	variantIDStr := r.PathValue("variant_id")

	variantID, err := strconv.ParseInt(variantIDStr, 10, 64)
	if err != nil || variantID <= 0 {
		apperror.Write(
			w,
			apperror.BadRequest(
				"INVALID_VARIANT_ID",
				"invalid variant id",
			),
		)
		return
	}
	result, err := h.inventoryService.GetByVariantID(r.Context(), variantID)
	if err != nil {
		apperror.Write(w, err)
		return
	}
	writeJSON(w, http.StatusOK, result)
}

// GetAll godoc
// @Summary Get all inventory
// @Description Retrieves all inventory records.
// @Tags Inventory
// @Produce json
// @Success 200 {array} dto.InventoryResponse
// @Failure 500 {object} map[string]any
// @Router /api/inventory [get]
func (h *InventoryHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	result, err := h.inventoryService.GetAll(r.Context())
	if err != nil {
		apperror.Write(w, err)
		return
	}

	writeJSON(w, http.StatusOK, result)

}

// Update godoc
// @Summary Update inventory
// @Description Updates inventory for a product variant.
// @Tags Inventory
// @Accept json
// @Produce json
// @Param variant_id path int true "Product Variant ID"
// @Param request body dto.UpdateInventoryRequest true "Inventory update data"
// @Security BearerAuth
// @Success 200 {object} dto.InventoryResponse
// @Failure 400 {object} map[string]any
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Failure 404 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /api/inventory/{variant_id} [put]
func (h *InventoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	variantIDStr := r.PathValue("variant_id")

	variantID, err := strconv.ParseInt(variantIDStr, 10, 64)
	if err != nil || variantID <= 0 {
		apperror.Write(
			w,
			apperror.BadRequest(
				"INVALID_VARIANT_ID",
				"invalid variant id",
			),
		)
		return
	}

	var req dto.UpdateInventoryRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.Write(
			w,
			apperror.BadRequest(
				"INVALID_REQUEST_BODY",
				"invalid request body",
			),
		)
		return
	}
	result, err := h.inventoryService.Update(
		r.Context(),
		variantID,
		&req,
	)
	if err != nil {
		apperror.Write(w, err)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

// Delete godoc
// @Summary Delete inventory
// @Description Deletes inventory for a product variant.
// @Tags Inventory
// @Produce json
// @Param variant_id path int true "Product Variant ID"
// @Security BearerAuth
// @Success 204
// @Failure 400 {object} map[string]any
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Failure 404 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /api/inventory/{variant_id} [delete]
func (h *InventoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	variantIDStr := r.PathValue("variant_id")

	variantID, err := strconv.ParseInt(variantIDStr, 10, 64)
	if err != nil || variantID < 0 {
		apperror.Write(
			w,
			apperror.BadRequest(
				"INVALID_VARIANT_ID",
				"invalid variant id",
			),
		)
		return
	}

	if err := h.inventoryService.Delete(r.Context(), variantID); err != nil {
		apperror.Write(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
