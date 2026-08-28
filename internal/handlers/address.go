package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ariiiiph/ecommerce/internal/apperror"
	"github.com/ariiiiph/ecommerce/internal/dto"
	"github.com/ariiiiph/ecommerce/internal/services"
)

type AddressHandler struct {
	addressService *services.AddressService
}

func NewAddressHandler(addressService *services.AddressService) *AddressHandler {
	return &AddressHandler{
		addressService: addressService,
	}
}

// Create godoc
// @Summary Create a new address
// @Description Creates a new address for a user.
// @Tags Addresses
// @Accept json
// @Produce json
// @Param request body dto.CreateAddressRequest true "Address data"
// @Security BearerAuth
// @Success 201 {object} dto.AddressResponse
// @Failure 400 {object} map[string]any
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Failure 404 {object} map[string]any
// @Failure 409 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /api/addresses [post]
func (h *AddressHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateAddressRequest

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

	result, err := h.addressService.Create(r.Context(), &req)
	if err != nil {
		apperror.Write(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, result)
}

// GetByID godoc
// @Summary Get an address by ID
// @Description Retrieves an address by its ID.
// @Tags Addresses
// @Produce json
// @Param id path int true "Address ID"
// @Security BearerAuth
// @Success 200 {object} dto.AddressResponse
// @Failure 400 {object} map[string]any
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Failure 404 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /api/addresses/{id} [get]
func (h *AddressHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		apperror.Write(
			w,
			apperror.BadRequest(
				"INVALID_ADDRESS_ID",
				"invalid address id",
			),
		)
		return
	}

	result, err := h.addressService.GetByID(r.Context(), id)
	if err != nil {
		apperror.Write(w, err)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

// GetAllByUserID godoc
// @Summary Get all addresses for a user
// @Description Retrieves all addresses belonging to a specific user.
// @Tags Addresses
// @Produce json
// @Param user_id path int true "User ID"
// @Security BearerAuth
// @Success 200 {array} dto.AddressResponse
// @Failure 400 {object} map[string]any
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Failure 404 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /api/users/{user_id}/addresses [get]
func (h *AddressHandler) GetAllByUserID(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.PathValue("user_id")

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil || userID <= 0 {
		apperror.Write(
			w,
			apperror.BadRequest(
				"INVALID_USER_ID",
				"invalid user id",
			),
		)
		return
	}

	result, err := h.addressService.GetAllByUserID(r.Context(), userID)
	if err != nil {
		apperror.Write(w, err)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

// Update godoc
// @Summary Update an address
// @Description Updates an existing address.
// @Tags Addresses
// @Accept json
// @Produce json
// @Param id path int true "Address ID"
// @Param request body dto.UpdateAddressRequest true "Updated address data"
// @Security BearerAuth
// @Success 200 {object} dto.AddressResponse
// @Failure 400 {object} map[string]any
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Failure 404 {object} map[string]any
// @Failure 409 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /api/addresses/{id} [put]
func (h *AddressHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		apperror.Write(
			w,
			apperror.BadRequest(
				"INVALID_ADDRESS_ID",
				"invalid address id",
			),
		)
		return
	}

	var req dto.UpdateAddressRequest

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

	result, err := h.addressService.Update(r.Context(), id, &req)
	if err != nil {
		apperror.Write(w, err)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

// Delete godoc
// @Summary Delete an address
// @Description Deletes an address by its ID.
// @Tags Addresses
// @Produce json
// @Param id path int true "Address ID"
// @Security BearerAuth
// @Success 204
// @Failure 400 {object} map[string]any
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Failure 404 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /api/addresses/{id} [delete]
func (h *AddressHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		apperror.Write(
			w,
			apperror.BadRequest(
				"INVALID_ADDRESS_ID",
				"invalid address id",
			),
		)
		return
	}

	if err := h.addressService.Delete(r.Context(), id); err != nil {
		apperror.Write(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
