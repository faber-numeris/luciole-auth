package extensions

import (
	"time"

	"github.com/faber-numeris/luciole-auth/authn/internal/adapters/inbound/httpapi/gen"
	"github.com/faber-numeris/luciole-auth/authn/internal/adapters/outbound/postgres/gen"
	"github.com/faber-numeris/luciole-auth/authn/internal/domain"
)

func StringToULID(s string) api.ULID {
	return api.ULID(s)
}

func StringToOptString(s string) api.OptString {
	return api.OptString{Set: true, Value: s}
}

func OptStringToString(os api.OptString) (string, error) {
	if v, ok := os.Get(); ok {
		return v, nil
	}
	return "", nil
}

func StringToOptstring(s string) api.OptString {
	if s == "" {
		return api.OptString{}
	}
	return api.OptString{Set: true, Value: s}
}

func TimeToTimePtr(t time.Time) (*time.Time, error) {
	return &t, nil
}

func StringToByteSlice(s string) ([]byte, error) {
	return []byte(s), nil
}

func UserTypeToApiUserType(t domain.UserType) api.UserType {
	switch t {
	case domain.UserTypeUser:
		return api.UserTypeUSER
	case domain.UserTypeServiceAccount:
		return api.UserTypeSERVICEACCOUNT
	case domain.UserTypeDevice:
		return api.UserTypeDEVICE
	default:
		return api.UserType(t)
	}
}

func SQLCUserToUserProfile(user gen.User) *domain.UserProfile {
	if user.FirstName == "" && user.LastName == "" && user.Locale == "" && user.Timezone == "" {
		return nil
	}
	return &domain.UserProfile{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Locale:    user.Locale,
		Timezone:  user.Timezone,
	}
}

func UserProfileToSQLCUser(profile *domain.UserProfile) (firstName, lastName, locale, timezone string) {
	if profile == nil {
		return
	}
	return profile.FirstName, profile.LastName, profile.Locale, profile.Timezone
}
