package postgres

import (
	"fmt"
	"sync"
	"time"

	"github.com/faber-numeris/luciole-auth/authn/internal/infrastructure/config"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

var DBInstance *sqlx.DB

var once sync.Once

func Connect() *sqlx.DB {
	once.Do(func() {
		cfg, err := config.LoadConfig()
		if err != nil {
			panic(fmt.Errorf("failed to load configuration: %w", err))
		}

		dsn := fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			cfg.DBHost(),
			cfg.DBPort(),
			cfg.DBUser(),
			cfg.DBPassword(),
			cfg.DBName(),
			cfg.DBSSLMode(),
		)

		db, err := sqlx.Connect("pgx", dsn)
		if err != nil {
			panic(fmt.Errorf("failed to connect to db: %w", err))
		}

		db.SetMaxOpenConns(10)
		db.SetMaxIdleConns(2)
		db.SetConnMaxLifetime(time.Hour)
		db.SetConnMaxIdleTime(30 * time.Minute)

		DBInstance = db
	})

	return DBInstance
}
