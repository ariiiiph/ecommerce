package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ariiiiph/ecommerce/internal/apperror"
	"github.com/ariiiiph/ecommerce/internal/dto"
	"github.com/ariiiiph/ecommerce/internal/services"
)

type AttributeValueHandler struct {
	attributeValueService *services.AttributeValueService
}

func NewAttributeValueHandler(attributeValueService *services.AttributeValueService) *AttributeValueHandler {
	return &AttributeValueHandler{
		attributeValueService: attributeValueService,
	}
}

// Create godoc
// @Summary Create a new attribute value
// @Description Creates a new value for an attribute.
// @Tags AttributeValues
// @Accept json
// @Produce json
// @Param request body dto.CreateAttributeValueRequest true "Attribute value data"
// @Security BearerAuth
// @Success 201 {object} dto.AttributeValueResponse
// @Failure 400 {object} map[string]any
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Failure 404 {object} map[string]any
// @Failure 409 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /api/attribute-values [post]
func (h *AttributeValueHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateAttributeValueRequest

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

	result, err := h.attributeValueService.Create(r.Context(), &req)
	if err != nil {
		apperror.Write(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, result)
}

// GetByID godoc
// @Summary Get attribute value by ID
// @Description Retrieves an attribute value by its ID.
// @Tags AttributeValues
// @Produce json
// @Param id path int true "Attribute Value ID"
// @Success 200 {object} dto.AttributeValueResponse
// @Failure 400 {object} map[string]any
// @Failure 404 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /api/attribute-values/{id} [get]
func (h *AttributeValueHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		apperror.Write(
			w,
			apperror.BadRequest(
				"INVALID_ATTRIBUTE_VALUE_ID",
				"invalid attribute value id",
			),
		)
		return
	}

	result, err := h.attributeValueService.GetByID(r.Context(), id)
	if err != nil {
		apperror.Write(w, err)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

// GetAll godoc
// @Summary Get all attribute values
// @Description Retrieves all attribute values.
// @Tags AttributeValues
// @Produce json
// @Success 200 {array} dto.AttributeValueResponse
// @Failure 500 {object} map[string]any
// @Router /api/attribute-values [get]
func (h *AttributeValueHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	result, err := h.attributeValueService.GetAll(r.Context())
	if err != nil {
		apperror.Write(w, err)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

// Update godoc
// @Summary Update attribute value
// @Description Updates an existing attribute value.
// @Tags AttributeValues
// @Accept json
// @Produce json
// @Param id path int true "Attribute Value ID"
// @Param request body dto.UpdateAttributeValueRequest true "Updated attribute value data"
// @Security BearerAuth
// @Success 200 {object} dto.AttributeValueResponse
// @Failure 400 {object} map[string]any
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Failure 404 {object} map[string]any
// @Failure 409 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /api/attribute-values/{id} [put]
func (h *AttributeValueHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		apperror.Write(
			w,
			apperror.BadRequest(
				"INVALID_ATTRIBUTE_VALUE_ID",
				"invalid attribute value id",
			),
		)
		return
	}

	var req dto.UpdateAttributeValueRequest

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

	result, err := h.attributeValueService.Update(r.Context(), id, &req)
	if err != nil {
		apperror.Write(w, err)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

// Delete godoc
// @Summary Delete attribute value
// @Description Deletes an attribute value by its ID.
// @Tags AttributeValues
// @Produce json
// @Param id path int true "Attribute Value ID"
// @Security BearerAuth
// @Success 204
// @Failure 400 {object} map[string]any
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Failure 404 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /api/attribute-values/{id} [delete]
func (h *AttributeValueHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		apperror.Write(
			w,
			apperror.BadRequest(
				"INVALID_ATTRIBUTE_VALUE_ID",
				"invalid attribute value id",
			),
		)
		return
	}

	if err := h.attributeValueService.Delete(r.Context(), id); err != nil {
		apperror.Write(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
