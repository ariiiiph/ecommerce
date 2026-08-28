package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ariiiiph/ecommerce/internal/apperror"
	"github.com/ariiiiph/ecommerce/internal/dto"
	"github.com/ariiiiph/ecommerce/internal/services"
)

type VariantAttributeValueHandler struct {
	variantAttributeValueService *services.VariantAttributeValueService
}

func NewVariantAttributeValueHandler(
	variantAttributeValueService *services.VariantAttributeValueService,
) *VariantAttributeValueHandler {
	return &VariantAttributeValueHandler{
		variantAttributeValueService: variantAttributeValueService,
	}
}

// Create godoc
// @Summary Assign an attribute value to a product variant
// @Description Assigns an attribute value to a product variant.
// @Tags VariantAttributeValues
// @Accept json
// @Produce json
// @Param request body dto.CreateVariantAttributeValueRequest true "Variant attribute value data"
// @Security BearerAuth
// @Success 201 {object} dto.VariantAttributeValueResponse
// @Failure 400 {object} map[string]any
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Failure 404 {object} map[string]any
// @Failure 409 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /api/variant-attribute-values [post]
func (h *VariantAttributeValueHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateVariantAttributeValueRequest

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

	result, err := h.variantAttributeValueService.Create(
		r.Context(),
		&req,
	)
	if err != nil {
		apperror.Write(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, result)
}

// GetByVariantID godoc
// @Summary Get attribute values of a product variant
// @Description Retrieves all attribute values assigned to a product variant.
// @Tags VariantAttributeValues
// @Produce json
// @Param variant_id path int true "Product Variant ID"
// @Success 200 {array} dto.VariantAttributeValueResponse
// @Failure 400 {object} map[string]any
// @Failure 404 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /api/product-variants/{variant_id}/attribute-values [get]
func (h *VariantAttributeValueHandler) GetByVariantID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("variant_id")

	variantID, err := strconv.ParseInt(idStr, 10, 64)
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

	result, err := h.variantAttributeValueService.GetByVariantID(
		r.Context(),
		variantID,
	)
	if err != nil {
		apperror.Write(w, err)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

// Delete godoc
// @Summary Remove an attribute value from a product variant
// @Description Removes an attribute value assignment from a product variant.
// @Tags VariantAttributeValues
// @Produce json
// @Param variant_id path int true "Product Variant ID"
// @Param attribute_value_id path int true "Attribute Value ID"
// @Security BearerAuth
// @Success 204
// @Failure 400 {object} map[string]any
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Failure 404 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /api/product-variants/{variant_id}/attribute-values/{attribute_value_id} [delete]
func (h *VariantAttributeValueHandler) Delete(w http.ResponseWriter, r *http.Request) {
	variantIDStr := r.PathValue("variant_id")
	attributeValueIDStr := r.PathValue("attribute_value_id")

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

	attributeValueID, err := strconv.ParseInt(
		attributeValueIDStr,
		10,
		64,
	)
	if err != nil || attributeValueID <= 0 {
		apperror.Write(
			w,
			apperror.BadRequest(
				"INVALID_ATTRIBUTE_VALUE_ID",
				"invalid attribute value id",
			),
		)
		return
	}

	if err := h.variantAttributeValueService.Delete(
		r.Context(),
		variantID,
		attributeValueID,
	); err != nil {
		apperror.Write(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
