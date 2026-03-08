package httpapi

import (
	"errors"
	"net/http"

	"github.com/faber-numeris/luciole-auth/authn/internal/adapters/httpapi/gen"
)

// MapError maps domain errors to ogen error responses
func MapError(err error) (int, api.Error) {
	// TODO: implement actual mapping when domain errors are defined
	if errors.Is(err, errors.New("user not found")) {
		return http.StatusNotFound, api.Error{
			Message: "User not found",
		}
	}
	return http.StatusInternalServerError, api.Error{
		Message: "Internal server error",
	}
}
