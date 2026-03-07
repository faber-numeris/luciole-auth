package config

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
)

type IDatabaseConfig interface {
	DBHost() string
	DBPort() int
	DBUser() string
	DBPassword() string
	DBName() string
	DBSSLMode() string
}

var _ IDatabaseConfig = (*DatabaseConfig)(nil)

type DatabaseConfig struct {
	DBHost_     string `env:"DB_HOST,required"`
	DBPort_     int    `env:"DB_PORT" envDefault:"5432"`
	DBUser_     string `env:"DB_USER,required"`
	DBPassword_ string `env:"DB_PASSWORD,required"`
	DBDBName_   string `env:"DB_NAME,required"`
	DBSSLMode_  string `env:"DB_SSLMODE,required"`
}

func (d DatabaseConfig) DBHost() string {
	return d.DBHost_
}

func (d DatabaseConfig) DBPort() int {
	return d.DBPort_
}

func (d DatabaseConfig) DBUser() string {
	return d.DBUser_
}

func (d DatabaseConfig) DBPassword() string {
	return d.DBPassword_
}

func (d DatabaseConfig) DBName() string {
	return d.DBDBName_
}

func (d DatabaseConfig) DBSSLMode() string {
	return d.DBSSLMode_
}

// DB wraps a pgxpool.Pool connection
var (
	instOnce sync.Once
	dbInst   *DB
	dbErr    error
)

type DB struct {
	*pgxpool.Pool
}

func GetInstance(cfg IDatabaseConfig) (*DB, error) {
	instOnce.Do(func() {
		dbInst, dbErr = newConnection(cfg)
	})
	return dbInst, dbErr
}

func CloseDB() error {
	if dbInst != nil {
		dbInst.Pool.Close()
		dbInst = nil
	}
	return nil
}

func newConnection(cfg IDatabaseConfig) (*DB, error) {
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
