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

type OrderHandler struct {
	orderService *services.OrderService
}

func NewOrderHandler(orderService *services.OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
	}
}

// Create godoc
// @Summary Create a new order
// @Description Creates an order from the authenticated user's cart.
// @Tags Orders
// @Accept json
// @Produce json
// @Param request body dto.CreateOrderRequest true "Order data"
// @Security BearerAuth
// @Success 201 {object} dto.OrderResponse
// @Failure 400 {object} map[string]any
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {object} map[string]any
// @Failure 409 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /api/orders [post]
func (h *OrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateOrderRequest

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

	claims, ok := middleware.GetClaims(r.Context())
	if !ok || claims == nil {
		apperror.Write(
			w,
			apperror.Unauthorized(
				"UNAUTHORIZED",
				"unauthorized",
			),
		)
		return
	}

	result, err := h.orderService.Create(
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
// @Summary Get an order by ID
// @Description Retrieves an order by its ID for the authenticated user.
// @Tags Orders
// @Produce json
// @Param id path int true "Order ID"
// @Security BearerAuth
// @Success 200 {object} dto.OrderResponse
// @Failure 400 {object} map[string]any
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /api/orders/{id} [get]
func (h *OrderHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		apperror.Write(
			w,
			apperror.BadRequest(
				"INVALID_ORDER_ID",
				"invalid order id",
			),
		)
		return
	}

	claims, ok := middleware.GetClaims(r.Context())
	if !ok || claims == nil {
		apperror.Write(
			w,
			apperror.Unauthorized(
				"UNAUTHORIZED",
				"unauthorized",
			),
		)
		return
	}

	result, err := h.orderService.GetByID(
		r.Context(),
		id,
	)
	if err != nil {
		apperror.Write(w, err)
		return
	}

	// Make sure the user can only access their own order.
	if result.UserID != claims.UserID {
		apperror.Write(
			w,
			apperror.NotFound(
				"ORDER_NOT_FOUND",
				"order not found",
			),
		)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

// GetAll godoc
// @Summary Get user's orders
// @Description Retrieves all orders belonging to the authenticated user.
// @Tags Orders
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.OrderResponse
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {object} map[string]any
// @Router /api/orders [get]
func (h *OrderHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetClaims(r.Context())
	if !ok || claims == nil {
		apperror.Write(
			w,
			apperror.Unauthorized(
				"UNAUTHORIZED",
				"unauthorized",
			),
		)
		return
	}

	result, err := h.orderService.GetAllByUserID(
		r.Context(),
		claims.UserID,
	)
	if err != nil {
		apperror.Write(w, err)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

// UpdateStatus godoc
// @Summary Update order status
// @Description Updates the status of an order. Admin only.
// @Tags Orders
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Param request body dto.UpdateOrderStatusRequest true "Order status"
// @Security BearerAuth
// @Success 200 {object} dto.OrderResponse
// @Failure 400 {object} map[string]any
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Failure 404 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /api/orders/{id}/status [patch]
func (h *OrderHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		apperror.Write(
			w,
			apperror.BadRequest(
				"INVALID_ORDER_ID",
				"invalid order id",
			),
		)
		return
	}

	var req dto.UpdateOrderStatusRequest

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

	result, err := h.orderService.UpdateStatus(
		r.Context(),
		id,
		req.Status,
	)
	if err != nil {
		apperror.Write(w, err)
		return
	}

	writeJSON(w, http.StatusOK, result)
}
