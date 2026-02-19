package model

import (
	"time"

	"github.com/faber-numeris/luciole-auth/authn/persistence/sqlc"
)

// User represents the domain model for a user entity
// This struct isolates the sqlc structs from the API generated objects
type User struct {
	ID          string    `json:"id"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	Password    string    `json:"password"`
	FirstName   string    `json:"firstName"`
	LastName    string    `json:"lastName"`
	PhoneNumber string    `json:"phoneNumber"`
	Bio         string    `json:"bio"`
	AvatarURL   string    `json:"avatarUrl"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// UserResponse represents the user data returned to API clients
type UserResponse struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateUserRequest represents the data required to create a new user
type CreateUserRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// UpdateUserRequest represents the data required to update an existing user
type UpdateUserRequest struct {
	Username *string `json:"username,omitempty" validate:"omitempty,min=3,max=50"`
	Email    *string `json:"email,omitempty" validate:"omitempty,email"`
}

// FromSQLC converts a sqlc.User to domain model.User
func (u *User) FromSQLC(sqlcUser sqlc.User) {
	u.ID = sqlcUser.ID
	u.Username = sqlcUser.Username
	u.Email = sqlcUser.Email
	// Note: sqlcUser.PasswordHash is []byte, but we don't expose it in the domain model
	if sqlcUser.CreatedAt != nil {
		u.CreatedAt = *sqlcUser.CreatedAt
	}
	if sqlcUser.UpdatedAt != nil {
		u.UpdatedAt = *sqlcUser.UpdatedAt
	}
}

// ToResponse converts domain model.User to UserResponse for API responses
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
