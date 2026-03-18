package postgres

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/faber-numeris/luciole-auth/authn/internal/infrastructure/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

var DBPoolInstance DBPool

var once sync.Once

func Connect() DBPool {
	once.Do(func() {
		ctx := context.Background()

		var cfg config.IDatabaseConfig = config.LoadConfig()

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
			panic(fmt.Errorf("failed to parse db config: %w", err))
		}

		poolCfg.MaxConns = 10
		poolCfg.MinConns = 2
		poolCfg.MaxConnLifetime = time.Hour
		poolCfg.MaxConnIdleTime = 30 * time.Minute
		poolCfg.HealthCheckPeriod = time.Minute

		pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
		if err != nil {
			panic(fmt.Errorf("failed to connect to db: %w", err))
		}

		// Always verify connection at startup
		if err := pool.Ping(ctx); err != nil {
			pool.Close()
			panic(fmt.Errorf("failed to ping db: %w", err))
		}

		DBPoolInstance = pool
	})

	return DBPoolInstance
}
