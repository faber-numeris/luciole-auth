package extensions

import (
	"testing"

	"github.com/faber-numeris/luciole-auth/authn/internal/adapters/inbound/httpapi/gen"
	"github.com/faber-numeris/luciole-auth/authn/internal/adapters/outbound/postgres/gen"
	"github.com/faber-numeris/luciole-auth/authn/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestExtensions(t *testing.T) {
	t.Run("StringToULID", func(t *testing.T) {
		assert.Equal(t, api.ULID("123"), StringToULID("123"))
	})

	t.Run("StringToOptString", func(t *testing.T) {
		opt := StringToOptString("test")
		assert.True(t, opt.Set)
		assert.Equal(t, "test", opt.Value)
	})

	t.Run("OptStringToString", func(t *testing.T) {
		s, err := OptStringToString(api.OptString{Set: true, Value: "test"})
		assert.NoError(t, err)
		assert.Equal(t, "test", s)

		s, err = OptStringToString(api.OptString{Set: false})
		assert.NoError(t, err)
		assert.Equal(t, "", s)
	})

	t.Run("StringToOptstring", func(t *testing.T) {
		opt := StringToOptstring("test")
		assert.True(t, opt.Set)
		assert.Equal(t, "test", opt.Value)

		opt = StringToOptstring("")
		assert.False(t, opt.Set)
	})

	t.Run("UserTypeToApiUserType", func(t *testing.T) {
		assert.Equal(t, api.UserTypeUSER, UserTypeToApiUserType(domain.UserTypeUser))
		assert.Equal(t, api.UserTypeSERVICEACCOUNT, UserTypeToApiUserType(domain.UserTypeServiceAccount))
		assert.Equal(t, api.UserTypeDEVICE, UserTypeToApiUserType(domain.UserTypeDevice))
		assert.Equal(t, api.UserType("UNKNOWN"), UserTypeToApiUserType(domain.UserType("UNKNOWN")))
	})

	t.Run("SQLCUserToUserProfile", func(t *testing.T) {
		user := gen.User{FirstName: "First", LastName: "Last"}
		profile := SQLCUserToUserProfile(user)
		assert.NotNil(t, profile)
		assert.Equal(t, "First", profile.FirstName)

		userEmpty := gen.User{}
		profileEmpty := SQLCUserToUserProfile(userEmpty)
		assert.Nil(t, profileEmpty)
	})

	t.Run("UserProfileToSQLCUser", func(t *testing.T) {
		profile := &domain.UserProfile{FirstName: "First"}
		f, _, _, _ := UserProfileToSQLCUser(profile)
		assert.Equal(t, "First", f)

		f2, _, _, _ := UserProfileToSQLCUser(nil)
		assert.Empty(t, f2)
	})
}
