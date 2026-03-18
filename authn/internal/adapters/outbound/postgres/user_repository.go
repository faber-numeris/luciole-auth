package postgresadapter

import (
	"context"
	"database/sql"
	"errors"

	"github.com/faber-numeris/luciole-auth/authn/internal/domain"
	outboundport "github.com/faber-numeris/luciole-auth/authn/internal/ports/outbound"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
)

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) outboundport.UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) CreateUser(ctx context.Context, user *domain.User, passwordHash []byte) (*domain.User, error) {
	query := `INSERT INTO users (email, password_hash) VALUES (?, ?) RETURNING id, email, password_hash, first_name, last_name, locale, timezone, created_at, updated_at, deleted_at`
	query = r.db.Rebind(query)

	var row userRow
	err := r.db.GetContext(ctx, &row, query, user.Email, passwordHash)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
			return nil, domain.ErrUserAlreadyExists
		}
		return nil, err
	}

	return r.toDomain(row), nil
}

func (r *userRepository) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	query := `SELECT id, email, password_hash, first_name, last_name, locale, timezone, created_at, updated_at, deleted_at FROM users WHERE id = ? LIMIT 1`
	query = r.db.Rebind(query)

	var row userRow
	err := r.db.GetContext(ctx, &row, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return r.toDomain(row), nil
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `SELECT id, email, password_hash, first_name, last_name, locale, timezone, created_at, updated_at, deleted_at FROM users WHERE email = ? AND deleted_at IS NULL LIMIT 1`
	query = r.db.Rebind(query)

	var row userRow
	err := r.db.GetContext(ctx, &row, query, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return r.toDomain(row), nil
}

func (r *userRepository) GetUserCredentials(ctx context.Context, email string) (*domain.UserCredentials, error) {
	query := `SELECT email, password_hash FROM users WHERE email = ? AND deleted_at IS NULL LIMIT 1`
	query = r.db.Rebind(query)

	var row struct {
		Email        string `db:"email"`
		PasswordHash []byte `db:"password_hash"`
	}
	err := r.db.GetContext(ctx, &row, query, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &domain.UserCredentials{
		Email:        row.Email,
		PasswordHash: row.PasswordHash,
	}, nil
}

func (r *userRepository) UpdateUser(ctx context.Context, user *domain.User) error {
	var firstName, lastName, locale, timezone string
	if user.Profile != nil {
		firstName = user.Profile.FirstName
		lastName = user.Profile.LastName
		locale = user.Profile.Locale
		timezone = user.Profile.Timezone
	}

	query := `UPDATE users SET email = ?, first_name = ?, last_name = ?, locale = ?, timezone = ?, updated_at = NOW() WHERE id = ?`
	query = r.db.Rebind(query)

	_, err := r.db.ExecContext(ctx, query, user.Email, firstName, lastName, locale, timezone, user.ID)
	return err
}

func (r *userRepository) DeleteUser(ctx context.Context, id string) error {
	query := `UPDATE users SET deleted_at = NOW() WHERE id = ?`
	query = r.db.Rebind(query)

	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *userRepository) ListUsers(ctx context.Context, params *outboundport.ListUsersParams) ([]*domain.User, error) {
	query := `
		SELECT id, email, password_hash, first_name, last_name, locale, timezone, created_at, updated_at, deleted_at
		FROM users
		WHERE (? IS NULL OR email = ?)
		  AND (?::TIMESTAMP IS NULL OR created_at >= ?)
		  AND (?::TIMESTAMP IS NULL OR created_at <= ?)
		  AND (deleted_at IS NULL) = ?
	`
	query = r.db.Rebind(query)

	var rows []userRow
	err := r.db.SelectContext(ctx, &rows, query,
		params.Email, params.Email,
		params.CreatedStartRange, params.CreatedStartRange,
		params.CreatedEndRange, params.CreatedEndRange,
		params.Active,
	)
	if err != nil {
		return nil, err
	}

	result := make([]*domain.User, len(rows))
	for i, row := range rows {
		result[i] = r.toDomain(row)
	}

	return result, nil
}

func (r *userRepository) UpdatePassword(ctx context.Context, userID string, passwordHash []byte) error {
	query := `UPDATE users SET password_hash = ?, updated_at = NOW() WHERE id = ?`
	query = r.db.Rebind(query)

	_, err := r.db.ExecContext(ctx, query, passwordHash, userID)
	return err
}

func (r *userRepository) toDomain(row userRow) *domain.User {
	return &domain.User{
		ID:    row.ID,
		Email: row.Email,
		Profile: &domain.UserProfile{
			FirstName: row.FirstName,
			LastName:  row.LastName,
			Locale:    row.Locale,
			Timezone:  row.Timezone,
		},
	}
}
