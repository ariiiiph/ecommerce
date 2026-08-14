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

func (h *CategoryHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	result, err := h.categoryService.GetAll(r.Context())
	if err != nil {
		http.Error(w, "failed to get categories", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

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
