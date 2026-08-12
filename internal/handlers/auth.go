package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/ariiiiph/ecommerce/internal/dto"
	"github.com/ariiiiph/ecommerce/internal/services"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register godoc
// @Summary Register a new user
// @Description Creates a new customer account.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Registration data"
// @Success 201 {object} dto.AuthResponse
// @Failure 400 {string} string "Invalid request body"
// @Failure 409 {string} string "Email already exists"
// @Failure 500 {string} string "Failed to register user"
// @Router /api/auth/register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" || req.Email == "" || req.Password == "" {
		http.Error(w, "name, email and password are required", http.StatusBadRequest)
		return
	}

	result, err := h.authService.Register(r.Context(), req)
	if err != nil {
		if errors.Is(err, services.ErrEmailAlreadyExists) {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}

		http.Error(w, "failed to register user", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, result)
}

// Login godoc
// @Summary Login user
// @Description Authenticates a user and returns access and refresh tokens.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login credentials"
// @Success 200 {object} dto.AuthResponse
// @Failure 400 {string} string "Invalid request body"
// @Failure 401 {string} string "Invalid email or password"
// @Failure 500 {string} string "Failed to login"
// @Router /api/auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		http.Error(w, "email and password are required", http.StatusBadRequest)
		return
	}

	result, err := h.authService.Login(r.Context(), req)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			http.Error(w, "invalid email or password", http.StatusUnauthorized)
			return
		}

		http.Error(w, "failed to login", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

// Refresh godoc
// @Summary Refresh access token
// @Description Generates a new access token and refresh token using a valid refresh token.
// @Tags Authentication
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.AuthResponse
// @Failure 401 {string} string "Invalid refresh token"
// @Failure 500 {string} string "Failed to refresh token"
// @Router /api/auth/refresh [post]
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	rawToken := extractBearerToken(r)

	if rawToken == "" {
		http.Error(w, "refresh token is required", http.StatusUnauthorized)
		return
	}

	result, err := h.authService.Refresh(r.Context(), rawToken)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			http.Error(w, "invalid refresh token", http.StatusUnauthorized)
			return
		}

		http.Error(w, "failed to refresh token", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

// Logout godoc
// @Summary Logout user
// @Description Revokes the current refresh token.
// @Tags Authentication
// @Security BearerAuth
// @Success 204
// @Failure 401 {string} string "Invalid refresh token"
// @Failure 500 {string} string "Failed to logout"
// @Router /api/auth/logout [post]
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	rawToken := extractBearerToken(r)

	if rawToken == "" {
		http.Error(w, "refresh token is required", http.StatusUnauthorized)
		return
	}

	err := h.authService.Logout(r.Context(), rawToken)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			http.Error(w, "invalid refresh token", http.StatusUnauthorized)
			return
		}

		http.Error(w, "failed to logout", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func extractBearerToken(r *http.Request) string {
	header := r.Header.Get("Authorization")

	if header == "" {
		return ""
	}

	parts := strings.SplitN(header, " ", 2)

	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}

	return parts[1]
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(data)
}
