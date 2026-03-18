package postgresadapter

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/faber-numeris/luciole-auth/authn/internal/domain"
	outboundport "github.com/faber-numeris/luciole-auth/authn/internal/ports/outbound"
	"github.com/jmoiron/sqlx"
)

type userConfirmationRepository struct {
	db *sqlx.DB
}

func NewUserConfirmationRepository(db *sqlx.DB) outboundport.UserConfirmationRepository {
	return &userConfirmationRepository{
		db: db,
	}
}

func (r *userConfirmationRepository) CreateUserConfirmation(ctx context.Context, userID string, token string, expiresAt time.Time) (*domain.UserConfirmation, error) {
	query := `INSERT INTO user_confirmations (user_id, token, expires_at) VALUES (?, ?, ?) RETURNING id, user_id, token, expires_at, confirmed_at, created_at, updated_at`
	query = r.db.Rebind(query)

	var row userConfirmationRow
	err := r.db.GetContext(ctx, &row, query, userID, token, expiresAt)
	if err != nil {
		return nil, err
	}

	return r.toDomain(row), nil
}

func (r *userConfirmationRepository) GetUserConfirmationByToken(ctx context.Context, token string) (string, error) {
	query := `SELECT user_id FROM user_confirmations WHERE token = ? AND expires_at > NOW() LIMIT 1`
	query = r.db.Rebind(query)

	var userID string
	err := r.db.GetContext(ctx, &userID, query, token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}
		return "", err
	}
	return userID, nil
}

func (r *userConfirmationRepository) ConfirmUserRegistration(ctx context.Context, userID string) error {
	query := `UPDATE user_confirmations SET confirmed_at = NOW(), updated_at = NOW() WHERE user_id = ? AND confirmed_at IS NULL`
	query = r.db.Rebind(query)

	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}

func (r *userConfirmationRepository) DeleteUserConfirmation(ctx context.Context, userID string) error {
	query := `DELETE FROM user_confirmations WHERE user_id = ?`
	query = r.db.Rebind(query)

	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}

func (r *userConfirmationRepository) toDomain(row userConfirmationRow) *domain.UserConfirmation {
	return &domain.UserConfirmation{
		ID:          row.ID,
		UserID:      row.UserID,
		Token:       row.Token,
		ExpiresAt:   row.ExpiresAt,
		ConfirmedAt: row.ConfirmedAt,
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
	}
}
