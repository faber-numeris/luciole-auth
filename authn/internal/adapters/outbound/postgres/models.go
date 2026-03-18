package postgresadapter

import (
	"time"
)

type userRow struct {
	ID           string     `db:"id"`
	Email        string     `db:"email"`
	PasswordHash []byte     `db:"password_hash"`
	FirstName    string     `db:"first_name"`
	LastName     string     `db:"last_name"`
	Locale       string     `db:"locale"`
	Timezone     string     `db:"timezone"`
	CreatedAt    time.Time  `db:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at"`
	DeletedAt    *time.Time `db:"deleted_at"`
}

type userConfirmationRow struct {
	ID          string     `db:"id"`
	UserID      string     `db:"user_id"`
	Token       string     `db:"token"`
	ExpiresAt   time.Time  `db:"expires_at"`
	ConfirmedAt *time.Time `db:"confirmed_at"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
}
