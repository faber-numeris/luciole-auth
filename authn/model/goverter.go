package model

import (
	"time"

	api "github.com/faber-numeris/luciole-auth/api/gen"
	"github.com/faber-numeris/luciole-auth/authn/persistence/sqlc"
	"github.com/google/uuid"
)

// StringToUUID converts string to uuid.UUID (kept for potential future use)
func StringToUUID(s string) (uuid.UUID, error) {
	return uuid.Parse(s)
}

// StringToULID converts string to api.ULID
func StringToULID(s string) (api.ULID, error) {
	return api.ULID(s), nil
}

// StringToOptString converts string to api.OptString
func StringToOptString(s string) (api.OptString, error) {
	return api.OptString{Set: true, Value: s}, nil
}

// TimeToTimePtr converts time.Time to *time.Time
func TimeToTimePtr(t time.Time) (*time.Time, error) {
	return &t, nil
}

// StringToByteSlice converts string to []byte
func StringToByteSlice(s string) ([]byte, error) {
	return []byte(s), nil
}

// goverter:converter
// goverter:output:file ./generated/goverter.gen.go
// goverter:ignoreMissing no
// goverter:skipCopySameType yes
// goverter:useZeroValueOnPointerInconsistency yes
// goverter:extend StringToULID
// goverter:extend StringToOptString
// goverter:extend TimeToTimePtr
// goverter:extend StringToByteSlice
type Converter interface {

	// goverter:ignore ID
	// goverter:ignore FirstName
	// goverter:ignore LastName
	// goverter:ignore PhoneNumber
	// goverter:ignore Bio
	// goverter:ignore AvatarURL
	// goverter:ignore CreatedAt
	// goverter:ignore UpdatedAt
	// goverter:ignore Password
	RegisterRequestToUserModel(src api.RegisterRequest) (*User, error)

	// Custom method to convert User to API UserResponse
	UserToRegisterUserRes(src User) (*api.UserResponse, error)

	// SQLC to Model conversions
	// goverter:ignore Password
	// goverter:ignore FirstName
	// goverter:ignore LastName
	// goverter:ignore PhoneNumber
	// goverter:ignore Bio
	// goverter:ignore AvatarURL
	// goverter:ignore UpdatedAt
	SQLCUserToUser(src sqlc.User) (*User, error)
}
