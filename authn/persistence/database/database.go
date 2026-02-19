package database

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/faber-numeris/luciole-auth/authn/configuration"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
)

var (
	instOnce sync.Once
	dbInst   *DB
	dbErr    error
)

type DB struct {
	*pgxpool.Pool
}

func GetInstance(cfg configuration.IDatabaseConfig) (*DB, error) {
	instOnce.Do(func() {
		dbInst, dbErr = newConnection(cfg)
	})
	return dbInst, dbErr
}

func Close() error {
	if dbInst != nil {
		dbInst.Pool.Close()
		dbInst = nil
	}
	return nil
}

func newConnection(cfg configuration.IDatabaseConfig) (*DB, error) {
	ctx := context.Background()

	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost(),
		cfg.DBPort(),
		cfg.DBUser(),
		cfg.DBPassword(),
		cfg.DBName(),
		cfg.DBSSLMode(),
	)

	poolCfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	poolCfg.MaxConns = 10
	poolCfg.MinConns = 2
	poolCfg.MaxConnLifetime = time.Hour
	poolCfg.MaxConnIdleTime = 30 * time.Minute
	poolCfg.HealthCheckPeriod = time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("create Pool: %w", err)
	}

	// Always verify connection at startup
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping db: %w", err)
	}

	return &DB{pool}, nil
}

func (db *DB) Close() error {
	db.Pool.Close()

	return nil
}

func (db *DB) Health(ctx context.Context) error {
	return db.Pool.Ping(ctx)
}
