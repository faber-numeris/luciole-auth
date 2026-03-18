package main

import (
	"log/slog"
	"os"
	"time"

	"github.com/faber-numeris/luciole-auth/authn/internal/app/bootstrap"
	"github.com/lmittmann/tint"
)

func main() {
	slog.SetDefault(slog.New(
		tint.NewHandler(os.Stderr, &tint.Options{
			Level:      slog.LevelDebug,
			TimeFormat: time.DateTime,
		}),
	))

	app := bootstrap.NewApp()
	if err := app.Run(); err != nil {
		slog.Error("Application error", "error", err)
	}
}
