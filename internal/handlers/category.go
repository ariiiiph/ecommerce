package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ariiiiph/ecommerce/internal/apperror"
	"github.com/ariiiiph/ecommerce/internal/dto"
	"github.com/ariiiiph/ecommerce/internal/services"
)

type CategoryHandler struct {
	categoryService *services.CategoryService
}

func NewCategoryHandler(categoryService *services.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
	}
}

// Create godoc
// @Summary Create a new category
// @Description Creates a new category.
// @Tags Categories
// @Accept json
// @Produce json
// @Param request body dto.CreateCategoryRequest true "Category data"
// @Security BearerAuth
// @Success 201 {object} dto.CategoryResponse
// @Failure 400 {object} map[string]any
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Failure 404 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /api/categories [post]
func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateCategoryRequest

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

	result, err := h.categoryService.Create(r.Context(), &req)
	if err != nil {
		apperror.Write(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, result)
}

// GetByID godoc
// @Summary Get category by ID
// @Description Retrieves a category by its ID.
// @Tags Categories
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} dto.CategoryResponse
// @Failure 400 {object} map[string]any
// @Failure 404 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /api/categories/{id} [get]
func (h *CategoryHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		apperror.Write(
			w,
			apperror.BadRequest(
				"INVALID_CATEGORY_ID",
				"invalid category id",
			),
		)
		return
	}

	result, err := h.categoryService.GetByID(r.Context(), id)
	if err != nil {
		apperror.Write(w, err)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

// GetAll godoc
// @Summary Get all categories
// @Description Retrieves all categories.
// @Tags Categories
// @Produce json
// @Success 200 {array} dto.CategoryResponse
// @Failure 500 {object} map[string]any
// @Router /api/categories [get]
func (h *CategoryHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	result, err := h.categoryService.GetAll(r.Context())
	if err != nil {
		apperror.Write(w, err)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

// Update godoc
// @Summary Update category
// @Description Updates an existing category.
// @Tags Categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Param request body dto.UpdateCategoryRequest true "Updated category data"
// @Security BearerAuth
// @Success 200 {object} dto.CategoryResponse
// @Failure 400 {object} map[string]any
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Failure 404 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /api/categories/{id} [put]
func (h *CategoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		apperror.Write(
			w,
			apperror.BadRequest(
				"INVALID_CATEGORY_ID",
				"invalid category id",
			),
		)
		return
	}

	var req dto.UpdateCategoryRequest

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

	result, err := h.categoryService.Update(r.Context(), id, &req)
	if err != nil {
		apperror.Write(w, err)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

// Delete godoc
// @Summary Delete category
// @Description Deletes a category by its ID.
// @Tags Categories
// @Produce json
// @Param id path int true "Category ID"
// @Security BearerAuth
// @Success 204
// @Failure 400 {object} map[string]any
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Failure 404 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /api/categories/{id} [delete]
func (h *CategoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		apperror.Write(
			w,
			apperror.BadRequest(
				"INVALID_CATEGORY_ID",
				"invalid category id",
			),
		)
		return
	}

	err = h.categoryService.Delete(r.Context(), id)
	if err != nil {
		apperror.Write(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetChildren godoc
// @Summary Get child categories
// @Description Retrieves all child categories of a given category.
// @Tags Categories
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {array} dto.CategoryResponse
// @Failure 400 {object} map[string]any
// @Failure 404 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /api/categories/{id}/children [get]
func (h *CategoryHandler) GetChildren(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		apperror.Write(
			w,
			apperror.BadRequest(
				"INVALID_CATEGORY_ID",
				"invalid category id",
			),
		)
		return
	}

	result, err := h.categoryService.GetChildren(r.Context(), id)
	if err != nil {
		apperror.Write(w, err)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

// GetTree godoc
// @Summary Get category tree
// @Description Retrieves the entire category tree starting from a given category.
// @Tags Categories
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {array} dto.CategoryResponse
// @Failure 400 {object} map[string]any
// @Failure 404 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /api/categories/{id}/tree [get]
func (h *CategoryHandler) GetTree(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		apperror.Write(
			w,
			apperror.BadRequest(
				"INVALID_CATEGORY_ID",
				"invalid category id",
			),
		)
		return
	}

	result, err := h.categoryService.GetTree(r.Context(), id)
	if err != nil {
		apperror.Write(w, err)
		return
	}

	writeJSON(w, http.StatusOK, result)
}
