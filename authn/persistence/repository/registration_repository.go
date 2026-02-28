package repository

import (
	"context"
	"time"

	"github.com/faber-numeris/luciole-auth/authn/persistence/sqlc"
)

type IRegistrationRepository interface {
	CreateRegistrationPending(ctx context.Context, email, code string, codeExpiresAt time.Time) (sqlc.RegistrationPending, error)
	GetRegistrationPending(ctx context.Context, id string) (sqlc.RegistrationPending, error)
	GetRegistrationPendingByEmail(ctx context.Context, email string) (sqlc.RegistrationPending, error)
	DeleteRegistrationPending(ctx context.Context, id string) error
	DeleteRegistrationPendingByEmail(ctx context.Context, email string) error
}

type RegistrationRepository struct {
	queries *sqlc.Queries
}

func NewRegistrationRepository(queries *sqlc.Queries) IRegistrationRepository {
	return &RegistrationRepository{
		queries: queries,
	}
}

func (r *RegistrationRepository) CreateRegistrationPending(ctx context.Context, email, code string, codeExpiresAt time.Time) (sqlc.RegistrationPending, error) {
	return r.queries.CreateRegistrationPending(ctx, sqlc.CreateRegistrationPendingParams{
		Email:         email,
		Code:          code,
		CodeExpiresAt: codeExpiresAt,
	})
}

func (r *RegistrationRepository) GetRegistrationPending(ctx context.Context, id string) (sqlc.RegistrationPending, error) {
	return r.queries.GetRegistrationPending(ctx, id)
}

func (r *RegistrationRepository) GetRegistrationPendingByEmail(ctx context.Context, email string) (sqlc.RegistrationPending, error) {
	return r.queries.GetRegistrationPendingByEmail(ctx, email)
}

func (r *RegistrationRepository) DeleteRegistrationPending(ctx context.Context, id string) error {
	return r.queries.DeleteRegistrationPending(ctx, id)
}

func (r *RegistrationRepository) DeleteRegistrationPendingByEmail(ctx context.Context, email string) error {
	return r.queries.DeleteRegistrationPendingByEmail(ctx, email)
}
