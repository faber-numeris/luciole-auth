package model

import (
	api "github.com/faber-numeris/luciole-auth/authn/api/gen"
	"github.com/faber-numeris/luciole-auth/authn/persistence/sqlc"
	"github.com/google/uuid"
)

// StringToUUID converts string to uuid.UUID (kept for potential future use)
func StringToUUID(s string) (uuid.UUID, error) {
	return uuid.Parse(s)
}

// goverter:converter
// goverter:output:file ./generated/goverter.gen.go
// goverter:ignoreMissing no
// goverter:skipCopySameType yes
// goverter:useZeroValueOnPointerInconsistency yes
// goverter:extend ./extensions:.*
type Converter interface {

	/*=================================================
		Conversions from Model to API types
	 ==================================================*/

	// goverter:ignore Type ID Contacts
	// goverter:map . Profile
	UserModelFromUserRequest(userRequest api.UserCreateRequest) (User, error)

	/*=================================================
	  Conversions from API types to Model
	 ==================================================*/

	// goverter:autoMap Profile
	UserModelToApiUser(user User) (api.User, error)

	/*=================================================
	  Conversions from SQLC types to Model
	 ==================================================*/

	// goverter:ignore Type Contacts
	// goverter:map . Profile
	UserModelFromSQLC(user sqlc.User) (User, error)

	/*=================================================
	  Conversions from Model to SQLC types
	 ==================================================*/

	// goverter:ignore PasswordHash CreatedAt UpdatedAt DeletedAt
	// goverter:autoMap Profile
	UserModelToSQLC(user User) (sqlc.User, error)
}
