package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ariiiiph/ecommerce/internal/apperror"
	"github.com/ariiiiph/ecommerce/internal/dto"
	"github.com/ariiiiph/ecommerce/internal/services"
)

type ProductHandler struct {
	productService *services.ProductService
}

func NewProductHandler(productService *services.ProductService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}

// Create godoc
// @Summary Create a new product
// @Description Creates a new product.
// @Tags Products
// @Accept json
// @Produce json
// @Param request body dto.CreateProductRequest true "Product data"
// @Security BearerAuth
// @Success 201 {object} dto.ProductResponse
// @Failure 400 {object} map[string]any
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Failure 404 {object} map[string]any
// @Failure 409 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /api/products [post]
func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateProductRequest

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

	result, err := h.productService.Create(r.Context(), &req)
	if err != nil {
		apperror.Write(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, result)
}

// GetByID godoc
// @Summary Get product by ID
// @Description Retrieves a product by its ID.
// @Tags Products
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} dto.ProductResponse
// @Failure 400 {object} map[string]any
// @Failure 404 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /api/products/{id} [get]
func (h *ProductHandler) GetByID(w http.ResponseWriter, r *http.Request) {
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

	result, err := h.productService.GetByID(r.Context(), id)
	if err != nil {
		apperror.Write(w, err)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

// GetAll godoc
// @Summary Get all products
// @Description Retrieves a paginated list of products.
// @Tags Products
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of products per page" default(10)
// @Success 200 {object} map[string]any
// @Failure 400 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /api/products [get]
func (h *ProductHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	page := 1
	limit := 10

	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		parsedPage, err := strconv.Atoi(pageStr)
		if err != nil {
			apperror.Write(
				w,
				apperror.BadRequest(
					"INVALID_PAGE",
					"page must be a valid integer",
				),
			)
			return
		}

		page = parsedPage
	}

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err != nil {
			apperror.Write(
				w,
				apperror.BadRequest(
					"INVALID_LIMIT",
					"limit must be a valid integer",
				),
			)
			return
		}

		limit = parsedLimit
	}

	req := &dto.PaginationRequest{
		Page:  page,
		Limit: limit,
	}

	result, pagination, err := h.productService.GetAll(r.Context(), req)
	if err != nil {
		apperror.Write(w, err)
		return
	}

	response := map[string]any{
		"data":       result,
		"pagination": pagination,
	}

	writeJSON(w, http.StatusOK, response)
}

// Update godoc
// @Summary Update product
// @Description Updates an existing product.
// @Tags Products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param request body dto.UpdateProductRequest true "Updated product data"
// @Security BearerAuth
// @Success 200 {object} dto.ProductResponse
// @Failure 400 {object} map[string]any
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Failure 404 {object} map[string]any
// @Failure 409 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /api/products/{id} [put]
func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
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

	var req dto.UpdateProductRequest

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

	result, err := h.productService.Update(r.Context(), id, &req)
	if err != nil {
		apperror.Write(w, err)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

// Delete godoc
// @Summary Delete product
// @Description Soft deletes a product by its ID.
// @Tags Products
// @Produce json
// @Param id path int true "Product ID"
// @Security BearerAuth
// @Success 204
// @Failure 400 {object} map[string]any
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Failure 404 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /api/products/{id} [delete]
func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
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

	if err := h.productService.Delete(r.Context(), id); err != nil {
		apperror.Write(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
