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

type WishlistItemHandler struct {
	wishlistItemService *services.WishlistItemService
}

func NewWishlistItemHandler(wishlistItemService *services.WishlistItemService) *WishlistItemHandler {
	return &WishlistItemHandler{
		wishlistItemService: wishlistItemService,
	}
}

// Create godoc
// @Summary Add product to wishlist
// @Description Adds a product to the authenticated user's wishlist.
// @Tags WishlistItems
// @Accept json
// @Produce json
// @Param request body dto.CreateWishlistItemRequest true "Wishlist item data"
// @Security BearerAuth
// @Success 201 {object} dto.WishlistItemResponse
// @Failure 400 {object} map[string]any
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {object} map[string]any
// @Failure 409 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /api/wishlist-items [post]
func (h *WishlistItemHandler) Create(w http.ResponseWriter, r *http.Request) {
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

	var req dto.CreateWishlistItemRequest

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

	result, err := h.wishlistItemService.Create(
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
// @Summary Get wishlist item by ID
// @Description Retrieves a wishlist item belonging to the authenticated user.
// @Tags WishlistItems
// @Produce json
// @Param id path int true "Wishlist Item ID"
// @Security BearerAuth
// @Success 200 {object} dto.WishlistItemResponse
// @Failure 400 {object} map[string]any
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {object} map[string]any
// @Failure 404 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /api/wishlist-items/{id} [get]
func (h *WishlistItemHandler) GetByID(w http.ResponseWriter, r *http.Request) {
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
				"INVALID_WISHLIST_ITEM_ID",
				"invalid wishlist item id",
			),
		)
		return
	}

	result, err := h.wishlistItemService.GetByID(
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

// GetAllByWishlistID godoc
// @Summary Get all wishlist items
// @Description Retrieves all items belonging to the authenticated user's wishlist.
// @Tags WishlistItems
// @Produce json
// @Param wishlist_id path int true "Wishlist ID"
// @Security BearerAuth
// @Success 200 {array} dto.WishlistItemResponse
// @Failure 400 {object} map[string]any
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {object} map[string]any
// @Failure 404 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /api/wishlists/{wishlist_id}/items [get]
func (h *WishlistItemHandler) GetAllByWishlistID(w http.ResponseWriter, r *http.Request) {
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

	idStr := r.PathValue("wishlist_id")

	wishlistID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || wishlistID <= 0 {
		apperror.Write(
			w,
			apperror.BadRequest(
				"INVALID_WISHLIST_ID",
				"invalid wishlist id",
			),
		)
		return
	}

	result, err := h.wishlistItemService.GetAllByWishlistID(
		r.Context(),
		wishlistID,
		claims.UserID,
	)
	if err != nil {
		apperror.Write(w, err)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

// Delete godoc
// @Summary Remove product from wishlist
// @Description Removes a wishlist item belonging to the authenticated user.
// @Tags WishlistItems
// @Produce json
// @Param id path int true "Wishlist Item ID"
// @Security BearerAuth
// @Success 204
// @Failure 400 {object} map[string]any
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {object} map[string]any
// @Failure 404 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /api/wishlist-items/{id} [delete]
func (h *WishlistItemHandler) Delete(w http.ResponseWriter, r *http.Request) {
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
				"INVALID_WISHLIST_ITEM_ID",
				"invalid wishlist item id",
			),
		)
		return
	}

	if err := h.wishlistItemService.Delete(
		r.Context(),
		id,
		claims.UserID,
	); err != nil {
		apperror.Write(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
