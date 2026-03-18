package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// DBPool defines an interface for a database connection pool, abstracting
// *pgxpool.Pool. This allows for easier mocking and dependency injection.
type DBPool interface {
	// Exec executes a command tag, like INSERT or UPDATE.
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	// Query executes a query that returns rows.
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	// QueryRow executes a query that is expected to return at most one row.
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	// Ping verifies a connection to the database is still alive.
	Ping(ctx context.Context) error
	// Close closes all connections in the pool.
	Close()
}