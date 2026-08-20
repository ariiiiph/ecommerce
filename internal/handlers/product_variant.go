package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ariiiiph/ecommerce/internal/apperror"
	"github.com/ariiiiph/ecommerce/internal/dto"
	"github.com/ariiiiph/ecommerce/internal/services"
)

type ProductVariantHandler struct {
	productVariantService *services.ProductVariantService
}

func NewProductVariantHandler(productVariantService *services.ProductVariantService) *ProductVariantHandler {
	return &ProductVariantHandler{
		productVariantService: productVariantService,
	}
}

// Create godoc
// @Summary Create a new product variant
// @Description Creates a new product variant
// @Tags ProductVariants
// @Accept json
// Produce json
// @Param request body dto.CreateProductVariantRequest true "Product variant data"
// @Security BearerAuth
// @Success 201 {object} dto.ProductVariantResponse
// @Failure 400 {object} map[string]any
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Failure 404 {object} map[string]any
// @Failure 409 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /api/product-variants [post]
func (h *ProductVariantHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateProductVariantRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.Write(
			w,
			apperror.BadRequest(
				"INVALID_REQUEST_BODY",
				"invalid body request",
			),
		)
		return
	}
	result, err := h.productVariantService.Create(r.Context(), &req)
	if err != nil {
		apperror.Write(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, result)
}

// GetByID godoc
// @Summary Get a product variant by ID
// @Description retrieves a product variant by its ID
// @Tags ProductVariants
// @Produce json
// @Param id path int true "Product Variant ID"
// @Success 200 {object} dto.ProductVariantResponse
// @Failure 400 {object} map[string]any
// @Failure 404 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /api/product-variants/{id} [get]
func (h *ProductVariantHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		apperror.Write(
			w,
			apperror.BadRequest(
				"INVALID_PRODUCT_VARIANT_ID",
				"invalid product variant id",
			),
		)
		return
	}
	result, err := h.productVariantService.GetByID(r.Context(), id)
	if err != nil {
		apperror.Write(w, err)
		return
	}
	writeJSON(w, http.StatusOK, result)
}

// GetAllByProductID godoc
// @Summary Get all variants by product ID
// @Description Retrieves all product variants belonging to a specific product.
// @Tags ProductVariants
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} dto.ProductVariantResponse
// @Failure 400 {object} map[string]any
// @Failure 404 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /api/products/{id}/variants [get]
func (h *ProductVariantHandler) GetAllByProductID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		apperror.Write(
			w,
			apperror.BadRequest(
				"INVALID_PRODUCT_ID",
				"invalid product id",
			),
		)
		return
	}
	result, err := h.productVariantService.GetAllByProductID(r.Context(), id)
	if err != nil {
		apperror.Write(w, err)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

// Update godoc
// @Summary update  product variant
// @Description updates an existing product variant
// @Tags ProductVariants
// @Accept json
// @Produce json
// @Param id path int true "Product variant ID"
// @Param request body dto.UpdateProductVariantRequest true "Updated product variant data"
// @Security BearerAuth
// @Success 200 {object} dto.ProductVariantResponse
// @Failure 400 {object} map[string]any
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Failure 404 {object} map[string]any
// @Failure 409 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /api/product-variants/{id} [put]
func (h *ProductVariantHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		apperror.Write(
			w,
			apperror.BadRequest(
				"INVALID_PRODUCT_VARIANT_ID",
				"invalid product variant id",
			),
		)
		return
	}

	var req dto.UpdateProductVariantRequest

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

	result, err := h.productVariantService.Update(r.Context(), id, &req)
	if err != nil {
		apperror.Write(w, err)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

// Delete godoc
// @Summary Delete product variant
// @Description deletes a product variant by its ID.
// @Tags ProductVariants
// @Produce json
// @Param id path int true "Product Variant ID"
// @Security BearerAuth
// @Success 204
// @Failure 400 {object} map[string]any
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Failure 404 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /api/product-variants/{id} [delete]
func (h *ProductVariantHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		apperror.Write(
			w,
			apperror.BadRequest(
				"INVALID_PRODUCT_VARIANT_ID",
				"invalid product variant id",
			),
		)
		return
	}
	err = h.productVariantService.Delete(r.Context(), id)
	if err != nil {
		apperror.Write(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
