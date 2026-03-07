package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/faber-numeris/luciole-auth/authn/internal/bootstrap/config"
)

type App struct {
	server *http.Server
}

func NewApp(cfg config.IAppConfig, handler http.Handler) *App {
	address := fmt.Sprintf(":%d", cfg.Port())
	return &App{
		server: &http.Server{
			Addr:    address,
			Handler: handler,
		},
	}
}

func (a *App) Run() error {
	slog.Info("Starting AuthN service", "address", a.server.Addr)

	srvErrChan := make(chan error, 1)
	go func() {
		if srvErr := a.server.ListenAndServe(); srvErr != nil && !errors.Is(srvErr, http.ErrServerClosed) {
			srvErrChan <- srvErr
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-srvErrChan:
		if err != nil {
			return fmt.Errorf("server error: %w", err)
		}
	case <-quit:
		slog.Info("Shutting down server...")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := a.server.Shutdown(ctx); err != nil {
			return fmt.Errorf("server forced to shutdown: %w", err)
		}

		if err := config.CloseDB(); err != nil {
			return fmt.Errorf("failed to close database: %w", err)
		}

		slog.Info("Server exited")
	}

	return nil
}
