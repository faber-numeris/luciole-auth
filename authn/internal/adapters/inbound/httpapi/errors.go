package httpapi

import (
	"errors"
	"net/http"

	"github.com/faber-numeris/luciole-auth/authn/internal/adapters/inbound/httpapi/gen"
	"github.com/faber-numeris/luciole-auth/authn/internal/domain"
)

// MapError maps domain errors to ogen error responses
func MapError(err error) (int, api.Error) {
	if errors.Is(err, domain.ErrUserAlreadyExists) {
		return http.StatusConflict, api.Error{
			Error:   "USER_ALREADY_EXISTS",
			Message: "User already exists",
		}
	}
	if errors.Is(err, domain.ErrUserNotFound) {
		return http.StatusNotFound, api.Error{
			Error:   "USER_NOT_FOUND",
			Message: "User not found",
		}
	}
	if errors.Is(err, domain.ErrInvalidRequest) {
		return http.StatusBadRequest, api.Error{
			Error:   "INVALID_REQUEST",
			Message: "Invalid request data",
		}
	}
	if errors.Is(err, domain.ErrInvalidToken) {
		return http.StatusBadRequest, api.Error{
			Error:   "INVALID_TOKEN",
			Message: "Invalid or missing token",
		}
	}
	if errors.Is(err, domain.ErrTokenExpired) {
		return http.StatusGone, api.Error{
			Error:   "TOKEN_EXPIRED",
			Message: "Token expired",
		}
	}
	if errors.Is(err, domain.ErrInvalidCredentials) {
		return http.StatusUnauthorized, api.Error{
			Error:   "INVALID_CREDENTIALS",
			Message: "Invalid credentials",
		}
	}
	if errors.Is(err, domain.ErrUnauthorized) {
		return http.StatusUnauthorized, api.Error{
			Error:   "UNAUTHORIZED",
			Message: "Unauthorized",
		}
	}
	if errors.Is(err, domain.ErrTooManyRequests) {
		return http.StatusTooManyRequests, api.Error{
			Error:   "TOO_MANY_REQUESTS",
			Message: "Too many requests",
		}
	}

	return http.StatusInternalServerError, api.Error{
		Error:   "INTERNAL_SERVER_ERROR",
		Message: "Internal server error",
	}
}
