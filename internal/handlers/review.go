package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ariiiiph/ecommerce/internal/apperror"
	"github.com/ariiiiph/ecommerce/internal/dto"
	"github.com/ariiiiph/ecommerce/internal/middleware"
	"github.com/ariiiiph/ecommerce/internal/services"
)

type ReviewHandler struct {
	reviewService *services.ReviewService
}

func NewReviewHandler(
	reviewService *services.ReviewService,
) *ReviewHandler {
	return &ReviewHandler{
		reviewService: reviewService,
	}
}

// Create godoc
// @Summary Create a review
// @Description Creates a new product review.
// @Tags Reviews
// @Accept json
// @Produce json
// @Param request body dto.CreateReviewRequest true "Review data"
// @Security BearerAuth
// @Success 201 {object} dto.ReviewResponse
// @Failure 400 {object} map[string]any
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {object} map[string]any
// @Failure 409 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /api/reviews [post]
func (h *ReviewHandler) Create(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetClaims(r.Context())
	if !ok {
		apperror.Write(
			w,
			apperror.Unauthorized(
				"UNAUTHORIZED",
				"unauthorized",
			),
		)
		return
	}

	var req dto.CreateReviewRequest

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

	result, err := h.reviewService.Create(
		r.Context(),
		claims.UserID,
		&req,
	)
	if err != nil {
		apperror.Write(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, result)
}

// GetByID godoc
// @Summary Get my review
// @Description Retrieves a review owned by the authenticated user.
// @Tags Reviews
// @Produce json
// @Param id path int true "Review ID"
// @Security BearerAuth
// @Success 200 {object} dto.ReviewResponse
// @Failure 400 {object} map[string]any
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /api/reviews/{id} [get]
func (h *ReviewHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetClaims(r.Context())
	if !ok {
		apperror.Write(
			w,
			apperror.Unauthorized(
				"UNAUTHORIZED",
				"unauthorized",
			),
		)
		return
	}

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		apperror.Write(
			w,
			apperror.BadRequest(
				"INVALID_REVIEW_ID",
				"invalid review id",
			),
		)
		return
	}

	result, err := h.reviewService.GetByID(
		r.Context(),
		id,
		claims.UserID,
	)
	if err != nil {
		apperror.Write(w, err)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

// GetAllByProductID godoc
// @Summary Get approved product reviews
// @Description Retrieves all approved reviews for a product.
// @Tags Reviews
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {array} dto.ReviewResponse
// @Failure 400 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /api/products/{id}/reviews [get]
func (h *ReviewHandler) GetAllByProductID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
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

	result, err := h.reviewService.GetAllByProductID(
		r.Context(),
		id,
	)
	if err != nil {
		apperror.Write(w, err)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

// GetAllByUserID godoc
// @Summary Get my reviews
// @Description Retrieves all reviews created by the authenticated user.
// @Tags Reviews
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.ReviewResponse
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {object} map[string]any
// @Router /api/reviews/me [get]
func (h *ReviewHandler) GetAllByUserID(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetClaims(r.Context())
	if !ok {
		apperror.Write(
			w,
			apperror.Unauthorized(
				"UNAUTHORIZED",
				"unauthorized",
			),
		)
		return
	}

	result, err := h.reviewService.GetAllByUserID(
		r.Context(),
		claims.UserID,
	)
	if err != nil {
		apperror.Write(w, err)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

// Update godoc
// @Summary Update review
// @Description Updates a review owned by the authenticated user.
// @Tags Reviews
// @Accept json
// @Produce json
// @Param id path int true "Review ID"
// @Param request body dto.UpdateReviewRequest true "Updated review data"
// @Security BearerAuth
// @Success 200 {object} dto.ReviewResponse
// @Failure 400 {object} map[string]any
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /api/reviews/{id} [put]
func (h *ReviewHandler) Update(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetClaims(r.Context())
	if !ok {
		apperror.Write(
			w,
			apperror.Unauthorized(
				"UNAUTHORIZED",
				"unauthorized",
			),
		)
		return
	}

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		apperror.Write(
			w,
			apperror.BadRequest(
				"INVALID_REVIEW_ID",
				"invalid review id",
			),
		)
		return
	}

	var req dto.UpdateReviewRequest

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

	result, err := h.reviewService.Update(
		r.Context(),
		id,
		claims.UserID,
		&req,
	)
	if err != nil {
		apperror.Write(w, err)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

// Delete godoc
// @Summary Delete review
// @Description Deletes a review owned by the authenticated user.
// @Tags Reviews
// @Produce json
// @Param id path int true "Review ID"
// @Security BearerAuth
// @Success 204
// @Failure 400 {object} map[string]any
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /api/reviews/{id} [delete]
func (h *ReviewHandler) Delete(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetClaims(r.Context())
	if !ok {
		apperror.Write(
			w,
			apperror.Unauthorized(
				"UNAUTHORIZED",
				"unauthorized",
			),
		)
		return
	}

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		apperror.Write(
			w,
			apperror.BadRequest(
				"INVALID_REVIEW_ID",
				"invalid review id",
			),
		)
		return
	}

	if err := h.reviewService.Delete(
		r.Context(),
		id,
		claims.UserID,
	); err != nil {
		apperror.Write(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetAllForAdmin godoc
// @Summary Get all reviews
// @Description Retrieves all reviews for admin moderation.
// @Tags Admin Reviews
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.ReviewResponse
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Failure 500 {object} map[string]any
// @Router /api/admin/reviews [get]
func (h *ReviewHandler) GetAllForAdmin(w http.ResponseWriter, r *http.Request) {
	result, err := h.reviewService.GetAllForAdmin(r.Context())
	if err != nil {
		apperror.Write(w, err)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

// GetAllByStatusForAdmin godoc
// @Summary Get reviews by status
// @Description Retrieves reviews filtered by status for admin moderation.
// @Tags Admin Reviews
// @Produce json
// @Param status query string true "Review status" Enums(pending,approved,rejected)
// @Security BearerAuth
// @Success 200 {array} dto.ReviewResponse
// @Failure 400 {object} map[string]any
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Failure 500 {object} map[string]any
// @Router /api/admin/reviews/status [get]
func (h *ReviewHandler) GetAllByStatusForAdmin(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")

	result, err := h.reviewService.GetAllByStatusForAdmin(
		r.Context(),
		status,
	)
	if err != nil {
		apperror.Write(w, err)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

// Approve godoc
// @Summary Approve a review
// @Description Approves a pending product review.
// @Tags Admin Reviews
// @Produce json
// @Param id path int true "Review ID"
// @Security BearerAuth
// @Success 204
// @Failure 400 {object} map[string]any
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Failure 404 {object} map[string]any
// @Failure 409 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /api/admin/reviews/{id}/approve [patch]
func (h *ReviewHandler) Approve(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		apperror.Write(
			w,
			apperror.BadRequest(
				"INVALID_REVIEW_ID",
				"invalid review id",
			),
		)
		return
	}

	if err := h.reviewService.Approve(r.Context(), id); err != nil {
		apperror.Write(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Reject godoc
// @Summary Reject a review
// @Description Rejects a pending product review.
// @Tags Admin Reviews
// @Produce json
// @Param id path int true "Review ID"
// @Security BearerAuth
// @Success 204
// @Failure 400 {object} map[string]any
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Failure 404 {object} map[string]any
// @Failure 409 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /api/admin/reviews/{id}/reject [patch]
func (h *ReviewHandler) Reject(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		apperror.Write(
			w,
			apperror.BadRequest(
				"INVALID_REVIEW_ID",
				"invalid review id",
			),
		)
		return
	}

	if err := h.reviewService.Reject(r.Context(), id); err != nil {
		apperror.Write(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
