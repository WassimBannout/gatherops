package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/WassimBannout/gatherops/internal/service"
	"github.com/go-chi/chi/v5/middleware"
)

type ErrorResponse struct {
	Error APIError `json:"error"`
}

type APIError struct {
	Code      string         `json:"code"`
	Message   string         `json:"message"`
	Details   map[string]any `json:"details,omitempty"`
	RequestID string         `json:"requestId,omitempty"`
}

func WriteJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func WriteError(w http.ResponseWriter, r *http.Request, status int, code, message string, details map[string]any) {
	WriteJSON(w, status, ErrorResponse{
		Error: APIError{
			Code:      code,
			Message:   message,
			Details:   details,
			RequestID: middleware.GetReqID(r.Context()),
		},
	})
}

func WriteServiceError(w http.ResponseWriter, r *http.Request, err error) {
	var serviceErr *service.Error
	if !errors.As(err, &serviceErr) {
		WriteError(w, r, http.StatusInternalServerError, "internal_error", "An unexpected error occurred", nil)
		return
	}

	switch serviceErr.Code {
	case service.ErrorCodeValidationFailed:
		WriteError(w, r, http.StatusBadRequest, string(serviceErr.Code), serviceErr.Message, serviceErr.Details)
	case service.ErrorCodeEmailAlreadyRegistered:
		WriteError(w, r, http.StatusConflict, string(serviceErr.Code), serviceErr.Message, serviceErr.Details)
	case service.ErrorCodeInvalidCredentials, service.ErrorCodeInvalidRefreshToken, service.ErrorCodeUnauthorized:
		WriteError(w, r, http.StatusUnauthorized, string(serviceErr.Code), serviceErr.Message, serviceErr.Details)
	case service.ErrorCodeForbidden:
		WriteError(w, r, http.StatusForbidden, string(serviceErr.Code), serviceErr.Message, serviceErr.Details)
	case service.ErrorCodeNotFound:
		WriteError(w, r, http.StatusNotFound, string(serviceErr.Code), serviceErr.Message, serviceErr.Details)
	default:
		WriteError(w, r, http.StatusInternalServerError, "internal_error", "An unexpected error occurred", nil)
	}
}
