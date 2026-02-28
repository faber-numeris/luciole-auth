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

	app "github.com/faber-numeris/luciole-auth/authn/app"
	"github.com/faber-numeris/luciole-auth/authn/di"
	"github.com/faber-numeris/luciole-auth/authn/persistence/database"
	"github.com/faber-numeris/luciole-auth/authn/tools"
	"github.com/lmittmann/tint"
)

func main() {

	conf := di.ProvideConfiguration()
	router := tools.Must(di.ProvideRouter())

	slog.SetDefault(slog.New(
		tint.NewHandler(os.Stderr, &tint.Options{
			Level:      slog.LevelDebug,
			TimeFormat: time.DateTime,
		}),
	))

	address := fmt.Sprintf(":%d", conf.Port())
	slog.Info("Starting AuthN service", "address", address)

	srv := app.NewServer(router)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	srvErrChan := make(chan error, 1)
	go func() {
		if err := srv.Run(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
			srvErrChan <- err
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-srvErrChan:
		if err != nil {
			slog.Error("Error starting server", "error", err)
		}
	case <-quit:
		slog.Info("Shutting down server...")

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			slog.Error("Server forced to shutdown", "error", err)
		}

		cancel()

		if err := database.Close(); err != nil {
			slog.Error("Failed to close database", "error", err)
		}

		slog.Info("Server exited")
	}

}
