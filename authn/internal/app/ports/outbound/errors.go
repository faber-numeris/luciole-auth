package outboundport

import "errors"

// Pre-defined errors for the outbound ports layer.
// These errors are used by adapters to communicate specific outcomes
// to the application core without leaking implementation details.
var (
	// ErrUserNotFound is returned by a UserRepository when a user cannot be found.
	ErrUserNotFound = errors.New("user not found")

	// ErrRepository is a generic error for unexpected issues within a repository adapter.
	ErrRepository = errors.New("repository error")
)
