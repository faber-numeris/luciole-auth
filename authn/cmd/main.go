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
	if err != nil {
		panic(err)
	}

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

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

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
