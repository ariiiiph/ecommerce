package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ariiiiph/ecommerce/internal/apperror"
	"github.com/ariiiiph/ecommerce/internal/dto"
	"github.com/ariiiiph/ecommerce/internal/services"
)

type CouponHandler struct {
	couponService *services.CouponService
}

func NewCouponHandler(couponService *services.CouponService) *CouponHandler {
	return &CouponHandler{
		couponService: couponService,
	}
}

// Create godoc
// @Summary Create a coupon
// @Description Creates a new coupon.
// @Tags Coupons
// @Accept json
// @Produce json
// @Param request body dto.CreateCouponRequest true "Coupon data"
// @Security BearerAuth
// @Success 201 {object} dto.CouponResponse
// @Failure 400 {object} map[string]any
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Failure 409 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /api/coupons [post]
func (h *CouponHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateCouponRequest

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

	result, err := h.couponService.Create(r.Context(), &req)
	if err != nil {
		apperror.Write(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, result)
}

// GetAll godoc
// @Summary Get all coupons
// @Description Retrieves all coupons.
// @Tags Coupons
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.CouponResponse
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Failure 500 {object} map[string]any
// @Router /api/coupons [get]
func (h *CouponHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	result, err := h.couponService.GetAll(r.Context())
	if err != nil {
		apperror.Write(w, err)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

// GetByID godoc
// @Summary Get coupon by ID
// @Description Retrieves a coupon by its ID.
// @Tags Coupons
// @Produce json
// @Param id path int true "Coupon ID"
// @Security BearerAuth
// @Success 200 {object} dto.CouponResponse
// @Failure 400 {object} map[string]any
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Failure 404 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /api/coupons/{id} [get]
func (h *CouponHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		apperror.Write(
			w,
			apperror.BadRequest(
				"INVALID_COUPON_ID",
				"invalid coupon id",
			),
		)
		return
	}

	result, err := h.couponService.GetByID(r.Context(), id)
	if err != nil {
		apperror.Write(w, err)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

// Update godoc
// @Summary Update coupon
// @Description Updates an existing coupon.
// @Tags Coupons
// @Accept json
// @Produce json
// @Param id path int true "Coupon ID"
// @Param request body dto.UpdateCouponRequest true "Updated coupon data"
// @Security BearerAuth
// @Success 200 {object} dto.CouponResponse
// @Failure 400 {object} map[string]any
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Failure 404 {object} map[string]any
// @Failure 409 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /api/coupons/{id} [put]
func (h *CouponHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		apperror.Write(
			w,
			apperror.BadRequest(
				"INVALID_COUPON_ID",
				"invalid coupon id",
			),
		)
		return
	}

	var req dto.UpdateCouponRequest

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

	result, err := h.couponService.Update(
		r.Context(),
		id,
		&req,
	)
	if err != nil {
		apperror.Write(w, err)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

// Delete godoc
// @Summary Delete coupon
// @Description Deletes a coupon by its ID.
// @Tags Coupons
// @Produce json
// @Param id path int true "Coupon ID"
// @Security BearerAuth
// @Success 204
// @Failure 400 {object} map[string]any
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Failure 404 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /api/coupons/{id} [delete]
func (h *CouponHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		apperror.Write(
			w,
			apperror.BadRequest(
				"INVALID_COUPON_ID",
				"invalid coupon id",
			),
		)
		return
	}

	if err := h.couponService.Delete(r.Context(), id); err != nil {
		apperror.Write(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
