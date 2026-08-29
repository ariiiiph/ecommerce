package handlers

import (
	"net/http"
	"strconv"

	"github.com/ariiiiph/ecommerce/internal/apperror"
	"github.com/ariiiiph/ecommerce/internal/middleware"
	"github.com/ariiiiph/ecommerce/internal/services"
)

type CartHandler struct {
	cartService *services.CartService
}

func NewCartHandler(cartService *services.CartService) *CartHandler {
	return &CartHandler{
		cartService: cartService,
	}
}

// Create godoc
// @Summary Create a new cart
// @Description Creates a new cart for the authenticated user.
// @Tags Carts
// @Produce json
// @Security BearerAuth
// @Success 201 {object} dto.CartResponse
// @Failure 400 {object} map[string]any
// @Failure 401 {string} string "Unauthorized"
// @Failure 409 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /api/carts [post]
func (h *CartHandler) Create(w http.ResponseWriter, r *http.Request) {
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

	result, err := h.cartService.Create(r.Context(), claims.UserID)
	if err != nil {
		apperror.Write(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, result)
}

// GetByID godoc
// @Summary Get a cart by ID
// @Description Retrieves a cart by its ID.
// @Tags Carts
// @Produce json
// @Param id path int true "Cart ID"
// @Security BearerAuth
// @Success 200 {object} dto.CartResponse
// @Failure 400 {object} map[string]any
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Failure 404 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /api/carts/{id} [get]
func (h *CartHandler) GetByID(w http.ResponseWriter, r *http.Request) {
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
	idStr := r.PathValue("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		apperror.Write(
			w,
			apperror.BadRequest(
				"INVALID_CART_ID",
				"invalid cart id",
			),
		)
		return
	}

	result, err := h.cartService.GetByID(r.Context(), id, claims.UserID)
	if err != nil {
		apperror.Write(w, err)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

// GetByUserID godoc
// @Summary Get the authenticated user's cart
// @Description Retrieves the cart belonging to the authenticated user.
// @Tags Carts
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.CartResponse
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /api/carts [get]
func (h *CartHandler) GetByUserID(w http.ResponseWriter, r *http.Request) {
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

	result, err := h.cartService.GetByUserID(r.Context(), claims.UserID)
	if err != nil {
		apperror.Write(w, err)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

// Delete godoc
// @Summary Delete a cart
// @Description Deletes a cart by its ID.
// @Tags Carts
// @Produce json
// @Param id path int true "Cart ID"
// @Security BearerAuth
// @Success 204
// @Failure 400 {object} map[string]any
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Failure 404 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /api/carts/{id} [delete]
func (h *CartHandler) Delete(w http.ResponseWriter, r *http.Request) {
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
	idStr := r.PathValue("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		apperror.Write(
			w,
			apperror.BadRequest(
				"INVALID_CART_ID",
				"invalid cart id",
			),
		)
		return
	}

	if err := h.cartService.Delete(r.Context(), id, claims.UserID); err != nil {
		apperror.Write(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
