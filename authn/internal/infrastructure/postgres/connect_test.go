package postgres

import (
	"context"
	"errors"
	"testing"

	"github.com/faber-numeris/luciole-auth/authn/internal/mocks"
	"github.com/agiledragon/gomonkey/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
)

func TestConnect(t *testing.T) {
	mockCfg := mocks.NewMockIDatabaseConfig(t)
	mockCfg.EXPECT().DBHost().Return("localhost").Maybe()
	mockCfg.EXPECT().DBPort().Return(5432).Maybe()
	mockCfg.EXPECT().DBUser().Return("user").Maybe()
	mockCfg.EXPECT().DBPassword().Return("pass").Maybe()
	mockCfg.EXPECT().DBName().Return("db").Maybe()
	mockCfg.EXPECT().DBSSLMode().Return("disable").Maybe()

	t.Run("success", func(t *testing.T) {
		dummyPool := &pgxpool.Pool{}
		
		patches := gomonkey.ApplyFunc(pgxpool.ParseConfig, func(connString string) (*pgxpool.Config, error) {
			return &pgxpool.Config{}, nil
		})
		defer patches.Reset()
		
		patches.ApplyFunc(pgxpool.NewWithConfig, func(ctx context.Context, c *pgxpool.Config) (*pgxpool.Pool, error) {
			return dummyPool, nil
		})
		
		patches.ApplyMethod(dummyPool, "Ping", func(_ *pgxpool.Pool, _ context.Context) error {
			return nil
		})

		pool := Connect(mockCfg)
		assert.Equal(t, dummyPool, pool)
	})

	t.Run("parse failure", func(t *testing.T) {
		patches := gomonkey.ApplyFunc(pgxpool.ParseConfig, func(connString string) (*pgxpool.Config, error) {
			return nil, errors.New("parse error")
		})
		defer patches.Reset()

		assert.Panics(t, func() {
			Connect(mockCfg)
		})
	})

	t.Run("connect failure", func(t *testing.T) {
		patches := gomonkey.ApplyFunc(pgxpool.ParseConfig, func(connString string) (*pgxpool.Config, error) {
			return &pgxpool.Config{}, nil
		})
		defer patches.Reset()
		
		patches.ApplyFunc(pgxpool.NewWithConfig, func(ctx context.Context, c *pgxpool.Config) (*pgxpool.Pool, error) {
			return nil, errors.New("connect error")
		})

		assert.Panics(t, func() {
			Connect(mockCfg)
		})
	})

	t.Run("ping failure", func(t *testing.T) {
		dummyPool := &pgxpool.Pool{}
		
		patches := gomonkey.ApplyFunc(pgxpool.ParseConfig, func(connString string) (*pgxpool.Config, error) {
			return &pgxpool.Config{}, nil
		})
		defer patches.Reset()
		
		patches.ApplyFunc(pgxpool.NewWithConfig, func(ctx context.Context, c *pgxpool.Config) (*pgxpool.Pool, error) {
			return dummyPool, nil
		})
		
		patches.ApplyMethod(dummyPool, "Ping", func(_ *pgxpool.Pool, _ context.Context) error {
			return errors.New("ping error")
		})
		
		patches.ApplyMethod(dummyPool, "Close", func(_ *pgxpool.Pool) {
			// mock close
		})

		assert.Panics(t, func() {
			Connect(mockCfg)
		})
	})
}
