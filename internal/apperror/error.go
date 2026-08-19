package apperror

import (
	"encoding/json"
	"net/http"
)

type AppError struct {
	Code    string
	Message string
	Status  int
}

func (e *AppError) Error() string {
	return e.Message
}

func NotFound(code, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Status:  404,
	}
}

func BadRequest(code, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Status:  400,
	}
}

func Unauthorized(code, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Status:  401,
	}
}

func Forbidden(code, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Status:  403,
	}
}

func Conflict(code, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Status:  409,
	}
}

func Internal(code, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Status:  500,
	}
}

func Write(w http.ResponseWriter, err error) {
	appErr, ok := err.(*AppError)

	if !ok {
		appErr = Internal(
			"INTERNAL_ERROR",
			"internal server error",
		)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(appErr.Status)

	_ = json.NewEncoder(w).Encode(map[string]any{
		"success": false,
		"error": map[string]string{
			"code":    appErr.Code,
			"message": appErr.Message,
		},
	})
}
