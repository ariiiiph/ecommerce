package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

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
// @Failure 400 {string} string "Invalid request body"
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Failure 404 {string} string "Parent category not found"
// @Failure 409 {string} string "Category already exists"
// @Failure 500 {string} string "Failed to create category"
// @Router /api/categories [post]
func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateCategoryRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	if req.Slug == "" {
		http.Error(w, "slug is required", http.StatusBadRequest)
		return
	}

	result, err := h.categoryService.Create(r.Context(), &req)
	if err != nil {
		if errors.Is(err, services.ErrCategoryAlreadyExists) {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}

		if errors.Is(err, services.ErrCategoryNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		http.Error(w, "failed to create category", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, result)
}

// GetByID godoc
// @Summary Get category by ID
// @Description Retrieves a category by its ID.
// @Tags Categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} dto.CategoryResponse
// @Failure 400 {string} string "Invalid category ID"
// @Failure 404 {string} string "Category not found"
// @Failure 500 {string} string "Failed to get category"
// @Router /api/categories/{id} [get]
func (h *CategoryHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		http.Error(w, "invalid category id", http.StatusBadRequest)
		return
	}

	result, err := h.categoryService.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, services.ErrCategoryNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		http.Error(w, "failed to get category", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

// GetAll godoc
// @Summary Get all categories
// @Description Retrieves all categories.
// @Tags Categories
// @Accept json
// @Produce json
// @Success 200 {array} dto.CategoryResponse
// @Failure 500 {string} string "Failed to get categories"
// @Router /api/categories [get]
func (h *CategoryHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	result, err := h.categoryService.GetAll(r.Context())
	if err != nil {
		http.Error(w, "failed to get categories", http.StatusInternalServerError)
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
// @Failure 400 {string} string "Invalid request body or category ID"
// @Failure 404 {string} string "Category not found"
// @Failure 500 {string} string "Failed to update category"
// @Router /api/categories/{id} [put]
func (h *CategoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		http.Error(w, "invalid category id", http.StatusBadRequest)
		return
	}

	var req dto.UpdateCategoryRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	if req.Slug == "" {
		http.Error(w, "slug is required", http.StatusBadRequest)
		return
	}

	result, err := h.categoryService.Update(r.Context(), id, &req)
	if err != nil {
		if errors.Is(err, services.ErrCategoryNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		http.Error(w, "failed to update category", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

// Delete godoc
// @Summary Delete category
// @Description Deletes a category by its ID.
// @Tags Categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Security BearerAuth
// @Success 204 {string} string "No Content"
// @Failure 400 {string} string "Invalid category ID"
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Failure 404 {string} string "Category not found"
// @Failure 500 {string} string "Failed to delete category"
// @Router /api/categories/{id} [delete]
func (h *CategoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		http.Error(w, "invalid category id", http.StatusBadRequest)
		return
	}

	err = h.categoryService.Delete(r.Context(), id)
	if err != nil {
		if errors.Is(err, services.ErrCategoryNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		http.Error(w, "failed to delete category", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetChildren godoc
// @Summary Get child categories
// @Description Retrieves all child categories of a given category.
// @Tags Categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {array} dto.CategoryResponse
// @Failure 400 {string} string "Invalid category ID"
// @Failure 404 {string} string "Category not found"
// @Failure 500 {string} string "Failed to get child categories"
// @Router /api/categories/{id}/children [get]
func (h *CategoryHandler) GetChildren(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		http.Error(w, "invalid category id", http.StatusBadRequest)
		return
	}

	result, err := h.categoryService.GetChildren(r.Context(), id)
	if err != nil {
		if errors.Is(err, services.ErrCategoryNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		http.Error(w, "failed to get child categories", http.StatusInternalServerError)
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
// @Failure 400 {string} string "Invalid category ID"
// @Failure 404 {string} string "Category not found"
// @Failure 500 {string} string "Failed to get category tree"
// @Router /api/categories/{id}/tree [get]
func (h *CategoryHandler) GetTree(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		http.Error(w, "invalid category id", http.StatusBadRequest)
		return
	}

	result, err := h.categoryService.GetTree(r.Context(), id)
	if err != nil {
		if errors.Is(err, services.ErrCategoryNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		http.Error(w, "failed to get category tree", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, result)
}
