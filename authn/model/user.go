package model

import (
	"time"
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
