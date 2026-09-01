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

type CartItemHandler struct {
	cartItemService *services.CartItemService
}

func NewCartItemHandler(cartItemService *services.CartItemService) *CartItemHandler {
	return &CartItemHandler{
		cartItemService: cartItemService,
	}
}

// Create godoc
// @Summary Add item to cart
// @Description Adds a product variant to the authenticated user's cart.
// @Tags CartItems
// @Accept json
// @Produce json
// @Param request body dto.CreateCartItemRequest true "Cart item data"
// @Security BearerAuth
// @Success 201 {object} dto.CartItemResponse
// @Failure 400 {object} map[string]any
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {object} map[string]any
// @Failure 409 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /api/cart-items [post]
func (h *CartItemHandler) Create(w http.ResponseWriter, r *http.Request) {
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

	var req dto.CreateCartItemRequest

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

	result, err := h.cartItemService.Create(
		r.Context(),
		&req,
		claims.UserID,
	)
	if err != nil {
		apperror.Write(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, result)
}

// GetByID godoc
// @Summary Get cart item by ID
// @Description Retrieves a cart item belonging to the authenticated user's cart.
// @Tags CartItems
// @Produce json
// @Param id path int true "Cart Item ID"
// @Security BearerAuth
// @Success 200 {object} dto.CartItemResponse
// @Failure 400 {object} map[string]any
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Failure 404 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /api/cart-items/{id} [get]
func (h *CartItemHandler) GetByID(w http.ResponseWriter, r *http.Request) {
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
				"INVALID_CART_ITEM_ID",
				"invalid cart item id",
			),
		)
		return
	}

	result, err := h.cartItemService.GetByID(
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

// GetAllByCartID godoc
// @Summary Get all cart items
// @Description Retrieves all items belonging to the authenticated user's cart.
// @Tags CartItems
// @Produce json
// @Param cart_id path int true "Cart ID"
// @Security BearerAuth
// @Success 200 {array} dto.CartItemResponse
// @Failure 400 {object} map[string]any
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Failure 404 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /api/carts/{cart_id}/items [get]
func (h *CartItemHandler) GetAllByCartID(w http.ResponseWriter, r *http.Request) {
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

	cartIDStr := r.PathValue("cart_id")

	cartID, err := strconv.ParseInt(cartIDStr, 10, 64)
	if err != nil || cartID <= 0 {
		apperror.Write(
			w,
			apperror.BadRequest(
				"INVALID_CART_ID",
				"invalid cart id",
			),
		)
		return
	}

	result, err := h.cartItemService.GetAllByCartID(
		r.Context(),
		cartID,
		claims.UserID,
	)
	if err != nil {
		apperror.Write(w, err)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

// Update godoc
// @Summary Update cart item
// @Description Updates the quantity of a cart item belonging to the authenticated user's cart.
// @Tags CartItems
// @Accept json
// @Produce json
// @Param id path int true "Cart Item ID"
// @Param request body dto.UpdateCartItemRequest true "Updated cart item data"
// @Security BearerAuth
// @Success 200 {object} dto.CartItemResponse
// @Failure 400 {object} map[string]any
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Failure 404 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /api/cart-items/{id} [put]
func (h *CartItemHandler) Update(w http.ResponseWriter, r *http.Request) {
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
				"INVALID_CART_ITEM_ID",
				"invalid cart item id",
			),
		)
		return
	}

	var req dto.UpdateCartItemRequest

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

	result, err := h.cartItemService.Update(
		r.Context(),
		id,
		&req,
		claims.UserID,
	)
	if err != nil {
		apperror.Write(w, err)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

// Delete godoc
// @Summary Delete cart item
// @Description Deletes a cart item belonging to the authenticated user's cart.
// @Tags CartItems
// @Produce json
// @Param id path int true "Cart Item ID"
// @Security BearerAuth
// @Success 204
// @Failure 400 {object} map[string]any
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Failure 404 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /api/cart-items/{id} [delete]
func (h *CartItemHandler) Delete(w http.ResponseWriter, r *http.Request) {
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
				"INVALID_CART_ITEM_ID",
				"invalid cart item id",
			),
		)
		return
	}

	if err := h.cartItemService.Delete(
		r.Context(),
		id,
		claims.UserID,
	); err != nil {
		apperror.Write(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
