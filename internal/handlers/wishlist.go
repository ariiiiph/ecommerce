package handlers

import (
	"net/http"
	"strconv"

	"github.com/ariiiiph/ecommerce/internal/apperror"
	"github.com/ariiiiph/ecommerce/internal/middleware"
	"github.com/ariiiiph/ecommerce/internal/services"
)

type WishlistHandler struct {
	wishlistService *services.WishlistService
}

func NewWishlistHandler(wishlistService *services.WishlistService) *WishlistHandler {
	return &WishlistHandler{
		wishlistService: wishlistService,
	}
}

// Create godoc
// @Summary Create a new wishlist
// @Description Creates a new wishlist for the authenticated user.
// @Tags Wishlists
// @Produce json
// @Security BearerAuth
// @Success 201 {object} dto.WishlistResponse
// @Failure 400 {object} map[string]any
// @Failure 401 {string} string "Unauthorized"
// @Failure 409 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /api/wishlists [post]
func (h *WishlistHandler) Create(w http.ResponseWriter, r *http.Request) {
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

	result, err := h.wishlistService.Create(r.Context(), claims.UserID)
	if err != nil {
		apperror.Write(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, result)
}

// GetByID godoc
// @Summary Get a wishlist by ID
// @Description Retrieves a wishlist by its ID.
// @Tags Wishlists
// @Produce json
// @Param id path int true "Wishlist ID"
// @Security BearerAuth
// @Success 200 {object} dto.WishlistResponse
// @Failure 400 {object} map[string]any
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Failure 404 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /api/wishlists/{id} [get]
func (h *WishlistHandler) GetByID(w http.ResponseWriter, r *http.Request) {
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
				"INVALID_WISHLIST_ID",
				"invalid wishlist id",
			),
		)
		return
	}

	result, err := h.wishlistService.GetByID(r.Context(), id, claims.UserID)
	if err != nil {
		apperror.Write(w, err)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

// GetByUserID godoc
// @Summary Get the authenticated user's wishlist
// @Description Retrieves the wishlist belonging to the authenticated user.
// @Tags Wishlists
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.WishlistResponse
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /api/wishlists [get]
func (h *WishlistHandler) GetByUserID(w http.ResponseWriter, r *http.Request) {
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

	result, err := h.wishlistService.GetByUserID(r.Context(), claims.UserID)
	if err != nil {
		apperror.Write(w, err)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

// Delete godoc
// @Summary Delete a wishlist
// @Description Deletes a wishlist by its ID.
// @Tags Wishlists
// @Produce json
// @Param id path int true "Wishlist ID"
// @Security BearerAuth
// @Success 204
// @Failure 400 {object} map[string]any
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Failure 404 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /api/wishlists/{id} [delete]
func (h *WishlistHandler) Delete(w http.ResponseWriter, r *http.Request) {
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
				"INVALID_WISHLIST_ID",
				"invalid wishlist id",
			),
		)
		return
	}

	if err := h.wishlistService.Delete(r.Context(), id, claims.UserID); err != nil {
		apperror.Write(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
