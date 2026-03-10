package main

import (
	"errors"
	"net/http"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/faber-numeris/luciole-auth/authn/internal/adapters/httpapi/gen"
	"github.com/faber-numeris/luciole-auth/authn/internal/infrastructure/config"
	infra_postgres "github.com/faber-numeris/luciole-auth/authn/internal/infrastructure/postgres"
	"github.com/faber-numeris/luciole-auth/authn/internal/mocks"
	"github.com/agiledragon/gomonkey/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewApp(t *testing.T) {
	mockCfg := mocks.NewMockIAppConfig(t)
	mockCfg.EXPECT().DBHost().Return("localhost").Maybe()
	mockCfg.EXPECT().Port().Return(8080).Maybe()
	mockCfg.EXPECT().AllowedOrigins().Return([]string{"*"}).Maybe()

	t.Run("success", func(t *testing.T) {
		patches := gomonkey.ApplyFunc(config.Load, func() (config.IAppConfig, error) {
			return mockCfg, nil
		})
		defer patches.Reset()

		patches.ApplyFunc(infra_postgres.Connect, func(_ config.IDatabaseConfig) *pgxpool.Pool {
			return &pgxpool.Pool{}
		})

		patches.ApplyFunc(api.NewServer, func(h api.Handler, s api.SecurityHandler, opts ...api.ServerOption) (*api.Server, error) {
			return &api.Server{}, nil
		})

		app := NewApp()
		assert.NotNil(t, app)
		assert.Equal(t, ":8080", app.server.Addr)
	})

	t.Run("config load failure", func(t *testing.T) {
		patches := gomonkey.ApplyFunc(config.Load, func() (config.IAppConfig, error) {
			return nil, errors.New("config error")
		})
		defer patches.Reset()

		assert.Panics(t, func() {
			NewApp()
		})
	})

	t.Run("server creation failure", func(t *testing.T) {
		patches := gomonkey.ApplyFunc(config.Load, func() (config.IAppConfig, error) {
			return mockCfg, nil
		})
		defer patches.Reset()

		patches.ApplyFunc(infra_postgres.Connect, func(_ config.IDatabaseConfig) *pgxpool.Pool {
			return &pgxpool.Pool{}
		})

		patches.ApplyFunc(api.NewServer, func(h api.Handler, s api.SecurityHandler, opts ...api.ServerOption) (*api.Server, error) {
			return nil, errors.New("server error")
		})

		assert.Panics(t, func() {
			NewApp()
		})
	})
}

type mockPool struct {
	mock.Mock
}

func (m *mockPool) Close() {
	m.Called()
}

func TestApp_Run(t *testing.T) {
	t.Run("shutdown on signal", func(t *testing.T) {
		mPool := &mockPool{}
		mPool.On("Close").Return()

		srv := &http.Server{
			Addr: ":0", // auto port
			Handler: http.NewServeMux(),
		}

		app := &App{
			server: srv,
			pool:   mPool,
		}

		// Run in goroutine
		errChan := make(chan error, 1)
		go func() {
			errChan <- app.Run()
		}()

		// Give it a moment to start
		time.Sleep(100 * time.Millisecond)

		// Send SIGINT
		process, _ := os.FindProcess(os.Getpid())
		process.Signal(syscall.SIGINT)

		select {
		case err := <-errChan:
			assert.NoError(t, err)
		case <-time.After(2 * time.Second):
			t.Fatal("app did not shut down in time")
		}

		mPool.AssertExpectations(t)
	})
}
