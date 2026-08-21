package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ariiiiph/ecommerce/internal/apperror"
	"github.com/ariiiiph/ecommerce/internal/dto"
	"github.com/ariiiiph/ecommerce/internal/services"
)

type ProductImageHandler struct {
	productImageService *services.ProductImageService
}

func NewProductImageHandler(productImageService *services.ProductImageService) *ProductImageHandler {
	return &ProductImageHandler{
		productImageService: productImageService,
	}
}

// Create godoc
// @Summary Create a new product image
// @Description Creates a new product image
// @Tags ProductImages
// @Accept json
// @Produce json
// @Param request body dto.CreateProductImageRequest true "Product image data"
// @Security BearerAuth
// @Success 201 {object} dto.ProductImageResponse
// @Failure 400 {object} map[string]any
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Failure 404 {object} map[string]any
// @Failure 409 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /api/product-images [post]
func (h *ProductImageHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateProductImageRequest

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

	result, err := h.productImageService.Create(r.Context(), &req)
	if err != nil {
		apperror.Write(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, result)

}

// GetByID godoc
// @Summary Get a product image by ID
// @Description retrieves a product image by its ID
// @Tags ProductImages
// @Produce json
// @Param id path int true "Product Image ID"
// @Success 200 {object} dto.ProductImageResponse
// @Failure 400 {object} map[string]any
// @Failure 404 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /api/product-images/{id} [get]
func (h *ProductImageHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		apperror.Write(
			w,
			apperror.BadRequest(
				"INVALID_PRODUCT_IMAGE_ID",
				"invalid prodct image id",
			),
		)
		return
	}

	result, err := h.productImageService.GetByID(r.Context(), id)
	if err != nil {
		apperror.Write(w, err)
		return
	}
	writeJSON(w, http.StatusOK, result)
}

// GetAllByProductID godoc
// @Summary Get all images by product ID
// @Description Retrieves all product images belonging to a specific product.
// @Tags ProductImages
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {array} dto.ProductImageResponse
// @Failure 400 {object} map[string]any
// @Failure 404 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /api/products/{id}/images [get]
func (h *ProductImageHandler) GetAllByProductID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		apperror.Write(
			w,
			apperror.NotFound(
				"INVALID_PRODUCT_ID",
				"invalid prodct id",
			),
		)
		return
	}
	result, err := h.productImageService.GetAllByProductID(r.Context(), id)
	if err != nil {
		apperror.Write(w, err)
		return
	}
	writeJSON(w, http.StatusOK, result)
}

// Update godoc
// @Summary update  product image
// @Description updates an existing product image
// @Tags ProductImages
// @Accept json
// @Produce json
// @Param id path int true "Product image ID"
// @Param request body dto.UpdateProductImageRequest true "Updated product image data"
// @Security BearerAuth
// @Success 200 {object} dto.ProductImageResponse
// @Failure 400 {object} map[string]any
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Failure 404 {object} map[string]any
// @Failure 409 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /api/product-images/{id} [put]
func (h *ProductImageHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		apperror.Write(
			w,
			apperror.BadRequest(
				"INVALID_PRODUCT_IMAGE_ID",
				"invalid product image id",
			),
		)
		return
	}

	var req dto.UpdateProductImageRequest

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

	result, err := h.productImageService.Update(r.Context(), id, &req)
	if err != nil {
		apperror.Write(w, err)
		return
	}

	writeJSON(w, http.StatusOK, result)

}

// Delete godoc
// @Summary Delete product image
// @Description deletes a product image by its ID.
// @Tags ProductImages
// @Produce json
// @Param id path int true "Product Image ID"
// @Security BearerAuth
// @Success 204
// @Failure 400 {object} map[string]any
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Failure 404 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /api/product-images/{id} [delete]
func (h *ProductImageHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		apperror.Write(
			w,
			apperror.BadRequest(
				"INVALID_PRODUCT_IMAGE_ID",
				"invalid product image id",
			),
		)
		return
	}

	err = h.productImageService.Delete(r.Context(), id)
	if err != nil {
		apperror.Write(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
