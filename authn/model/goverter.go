package model

import api "github.com/faber-numeris/luciole-auth/api/gen"

// goverter:converter
// goverter:output:file ./generated/goverter.gen.go
// goverter:ignoreMissing no
// goverter:skipCopySameType yes
// goverter:useZeroValueOnPointerInconsistency yes
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
}
