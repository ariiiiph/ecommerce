package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ariiiiph/ecommerce/internal/apperror"
	"github.com/ariiiiph/ecommerce/internal/dto"
	"github.com/ariiiiph/ecommerce/internal/services"
)

type AttributeHandler struct {
	attributeService *services.AttributeService
}

func NewAttributeHandler(attributeService *services.AttributeService) *AttributeHandler {
	return &AttributeHandler{
		attributeService: attributeService,
	}
}

// Create godoc
// @Summary Create a new attribute
// @Description Creates a new product attribute.
// @Tags Attributes
// @Accept json
// @Produce json
// @Param request body dto.CreateAttributeRequest true "Attribute data"
// @Security BearerAuth
// @Success 201 {object} dto.AttributeResponse
// @Failure 400 {object} map[string]any
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Failure 409 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /api/attributes [post]
func (h *AttributeHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateAttributeRequest

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

	result, err := h.attributeService.Create(r.Context(), &req)
	if err != nil {
		apperror.Write(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, result)
}

// GetByID godoc
// @Summary Get attribute by ID
// @Description Retrieves an attribute by its ID.
// @Tags Attributes
// @Produce json
// @Param id path int true "Attribute ID"
// @Success 200 {object} dto.AttributeResponse
// @Failure 400 {object} map[string]any
// @Failure 404 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /api/attributes/{id} [get]
func (h *AttributeHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		apperror.Write(
			w,
			apperror.BadRequest(
				"INVALID_ATTRIBUTE_ID",
				"invalid attribute id",
			),
		)
		return
	}

	result, err := h.attributeService.GetByID(r.Context(), id)
	if err != nil {
		apperror.Write(w, err)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

// GetAll godoc
// @Summary Get all attributes
// @Description Retrieves all product attributes.
// @Tags Attributes
// @Produce json
// @Success 200 {array} dto.AttributeResponse
// @Failure 500 {object} map[string]any
// @Router /api/attributes [get]
func (h *AttributeHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	result, err := h.attributeService.GetAll(r.Context())
	if err != nil {
		apperror.Write(w, err)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

// Update godoc
// @Summary Update attribute
// @Description Updates an existing product attribute.
// @Tags Attributes
// @Accept json
// @Produce json
// @Param id path int true "Attribute ID"
// @Param request body dto.UpdateAttributeRequest true "Updated attribute data"
// @Security BearerAuth
// @Success 200 {object} dto.AttributeResponse
// @Failure 400 {object} map[string]any
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Failure 404 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /api/attributes/{id} [put]
func (h *AttributeHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		apperror.Write(
			w,
			apperror.BadRequest(
				"INVALID_ATTRIBUTE_ID",
				"invalid attribute id",
			),
		)
		return
	}

	var req dto.UpdateAttributeRequest

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

	result, err := h.attributeService.Update(r.Context(), id, &req)
	if err != nil {
		apperror.Write(w, err)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

// Delete godoc
// @Summary Delete attribute
// @Description Deletes an attribute by its ID.
// @Tags Attributes
// @Produce json
// @Param id path int true "Attribute ID"
// @Security BearerAuth
// @Success 204
// @Failure 400 {object} map[string]any
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Failure 404 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /api/attributes/{id} [delete]
func (h *AttributeHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		apperror.Write(
			w,
			apperror.BadRequest(
				"INVALID_ATTRIBUTE_ID",
				"invalid attribute id",
			),
		)
		return
	}

	if err := h.attributeService.Delete(r.Context(), id); err != nil {
		apperror.Write(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
