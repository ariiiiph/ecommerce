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

type PaymentHandler struct {
	paymentService *services.PaymentService
}

func NewPaymentHandler(paymentService *services.PaymentService) *PaymentHandler {
	return &PaymentHandler{
		paymentService: paymentService,
	}
}

// Create godoc
// @Summary Create a payment
// @Description Creates a mock payment for the authenticated user's pending order.
// @Tags Payments
// @Accept json
// @Produce json
// @Param request body dto.CreatePaymentRequest true "Payment data"
// @Security BearerAuth
// @Success 201 {object} dto.PaymentResponse
// @Failure 400 {object} map[string]any
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {object} map[string]any
// @Failure 409 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /api/payments [post]
func (h *PaymentHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreatePaymentRequest

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

	result, err := h.paymentService.Create(
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
// @Summary Get a payment by ID
// @Description Retrieves a payment belonging to the authenticated user.
// @Tags Payments
// @Produce json
// @Param id path int true "Payment ID"
// @Security BearerAuth
// @Success 200 {object} dto.PaymentResponse
// @Failure 400 {object} map[string]any
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /api/payments/{id} [get]
func (h *PaymentHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		apperror.Write(
			w,
			apperror.BadRequest(
				"INVALID_PAYMENT_ID",
				"invalid payment id",
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

	result, err := h.paymentService.GetByID(
		r.Context(),
		claims.UserID,
		id,
	)
	if err != nil {
		apperror.Write(w, err)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

// GetByOrderID godoc
// @Summary Get payment by order ID
// @Description Retrieves the payment belonging to the authenticated user's order.
// @Tags Payments
// @Produce json
// @Param order_id path int true "Order ID"
// @Security BearerAuth
// @Success 200 {object} dto.PaymentResponse
// @Failure 400 {object} map[string]any
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /api/orders/{order_id}/payment [get]
func (h *PaymentHandler) GetByOrderID(w http.ResponseWriter, r *http.Request) {
	orderIDStr := r.PathValue("order_id")

	orderID, err := strconv.ParseInt(orderIDStr, 10, 64)
	if err != nil || orderID <= 0 {
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

	result, err := h.paymentService.GetByOrderID(
		r.Context(),
		claims.UserID,
		orderID,
	)
	if err != nil {
		apperror.Write(w, err)
		return
	}

	writeJSON(w, http.StatusOK, result)
}
