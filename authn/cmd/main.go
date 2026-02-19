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

	"github.com/faber-numeris/luciole-auth/authn/di"
	"github.com/faber-numeris/luciole-auth/authn/persistence/database"
	"github.com/lmittmann/tint"
)

func main() {

	conf := di.ProvideConfiguration()
	router, err := di.ProvideRouter()
	// TODO: Create the Must function to deal with require or die constraints
	// assignees: rafaelsousa
	if err != nil {
		panic(err)
	}

	// Give log levels different colors.
	slog.SetDefault(slog.New(
		tint.NewHandler(os.Stderr, &tint.Options{
			Level:      slog.LevelDebug,
			TimeFormat: time.DateTime,
		}),
	))

	address := fmt.Sprintf(":%d", conf.Port())
	slog.Info("Starting AuthN service", "address", address)

	srv := &http.Server{
		Addr:    address,
		Handler: router,
	}

	// TODO: Put this into a separate component called app that we'll be called from main.go
	// assignees: rafaelsousa
	srvErrChan := make(chan error, 1)
	go func() {
		if srvErr := srv.ListenAndServe(); srvErr != nil && !errors.Is(srvErr, http.ErrServerClosed) {
			srvErrChan <- err
		}
	}()

	// Gracefully stops the server on CTRL + C actions.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-srvErrChan:
		if err != nil {
			slog.Error("Error starting server", "error", err)
		}
	case <-quit:
		slog.Info("Shutting down server...")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			slog.Error("Server forced to shutdown", "error", err)
		}

		if err := database.Close(); err != nil {
			slog.Error("Failed to close database", "error", err)
		}

		slog.Info("Server exited")
	}

}
