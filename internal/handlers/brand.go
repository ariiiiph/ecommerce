package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ariiiiph/ecommerce/internal/apperror"
	"github.com/ariiiiph/ecommerce/internal/dto"
	"github.com/ariiiiph/ecommerce/internal/services"
)

type BrandHandler struct {
	brandService *services.BrandService
}

func NewBrandHandler(brandService *services.BrandService) *BrandHandler {
	return &BrandHandler{
		brandService: brandService,
	}
}

// Create godoc
// @Summary Create a new brand
// @Description Creates a new brand.
// @Tags Brands
// @Accept json
// @Produce json
// @Param request body dto.CreateBrandRequest true "Brand data"
// @Security BearerAuth
// @Success 201 {object} dto.BrandResponse
// @Failure 400 {object} map[string]any "Bad request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Failure 409 {object} map[string]any "Brand already exists"
// @Failure 500 {object} map[string]any "Internal server error"
// @Router /api/brands [post]
func (h *BrandHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateBrandRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.Write(w, apperror.BadRequest(
			"INVALID_REQUEST_BODY",
			"invalid request body",
		))
		return
	}

	result, err := h.brandService.Create(r.Context(), &req)
	if err != nil {
		apperror.Write(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, result)
}

// GetByID godoc
// @Summary Get brand by ID
// @Description Retrieves a brand by its iD
// @Tags Brands
// @Accept json
// @Produce json
// @Param id path int true "Brand ID"
// @Success 200 {object} dto.BrandResponse
// @Failure 400 {object} map[string]any "Bad request"
// @Failure 404 {object} map[string]any "Brand not found"
// @Failure 500 {object} map[string]any "Internal server error"
// @Router /api/brands/{id} [get]
func (h *BrandHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		apperror.Write(w, apperror.BadRequest(
			"INVALID_BRAND_ID",
			"invalid brand id",
		))
		return
	}

	result, err := h.brandService.GetByID(r.Context(), id)
	if err != nil {
		apperror.Write(w, err)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

// GetAll godoc
// @Summary Get all brands
// @Description Retrieves all brands.
// @Tags Brands
// @Accept json
// @Produce json
// @Success 200 {array} dto.BrandResponse
// @Failure 500 {string} string "Failed to get brands"
// @Router /api/brands [get]
func (h *BrandHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	result, err := h.brandService.GetAll(r.Context())
	if err != nil {
		apperror.Write(w, err)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

// Update godoc
// @Summary Update brand
// @Description Updates an existing brand.
// @Tags Brands
// @Accept json
// @Produce json
// @Param id path int true "Brand ID"
// @Param request body dto.UpdateBrandRequest true "Updated brand data"
// @Security BearerAuth
// @Success 200 {object} dto.BrandResponse
// @Failure 400 {object} map[string]any "Bad request"
// @Failure 404 {object} map[string]any "Brand not found"
// @Failure 500 {object} map[string]any "Internal server error"
// @Router /api/brands/{id} [put]
func (h *BrandHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		apperror.Write(w, apperror.BadRequest(
			"INVALID_BRAND_ID",
			"invalid brand id",
		))
		return
	}

	var req dto.UpdateBrandRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.Write(w, apperror.BadRequest(
			"INVALID_REQUEST_BODY",
			"invalid request body",
		))
		return
	}

	result, err := h.brandService.Update(r.Context(), id, &req)
	if err != nil {
		apperror.Write(w, err)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

// Delete godoc
// @Summary Delete brand
// @Description Deletes a brand by its ID.
// @Tags Brands
// @Accept json
// @Produce json
// @Param id path int true "Brand ID"
// @Security BearerAuth
// @Success 204
// @Failure 400 {object} map[string]any "Bad request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Failure 404 {object} map[string]any "Brand not found"
// @Failure 500 {object} map[string]any "Internal server error"
// @Router /api/brands/{id} [delete]
func (h *BrandHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		apperror.Write(w, apperror.BadRequest(
			"INVALID_BRAND_ID",
			"invalid brand id",
		))
		return
	}

	err = h.brandService.Delete(r.Context(), id)
	if err != nil {
		apperror.Write(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
