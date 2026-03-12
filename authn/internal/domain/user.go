package domain

import "time"

type UserType string

const (
	UserTypeUser           UserType = "USER"
	UserTypeServiceAccount UserType = "SERVICE_ACCOUNT"
	UserTypeDevice         UserType = "DEVICE"
)

type ContactType string

const (
	ContactTypeAddress   ContactType = "ADDRESS"
	ContactTypeTelephone ContactType = "TELEPHONE"
	ContactTypeMobile    ContactType = "MOBILE"
	ContactTypeEmail     ContactType = "EMAIL"
	ContactTypeWebsite   ContactType = "WEBSITE"
)

type Contact struct {
	ContactType  ContactType `json:"contactType"`
	ContactValue string      `json:"contactValue"`
}

type User struct {
	ID       string   `json:"id"`
	Type     UserType `json:"type"`
	Email    string   `json:"email"`
	Profile  *UserProfile
	Contacts []Contact
}

type UserProfile struct {
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
	Locale    string `json:"locale,omitempty"`
	Timezone  string `json:"timezone,omitempty"`
}

type UserCredentials struct {
	Email        string
	PasswordHash []byte
}

type ConfirmationToken struct {
	Token     string
	ExpiresAt time.Time
	UserID    string
}

type PasswordResetToken struct {
	Token     string
	ExpiresAt time.Time
	UserID    string
}

type UserConfirmation struct {
	ID          string
	UserID      string
	UserEmail   string
	Token       string
	ExpiresAt   time.Time
	ConfirmedAt *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
