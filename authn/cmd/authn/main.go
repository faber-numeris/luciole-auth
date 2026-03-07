package main

import (
	"log/slog"
	"os"
	"time"

	"github.com/faber-numeris/luciole-auth/authn/internal/bootstrap"
	"github.com/faber-numeris/luciole-auth/authn/internal/platform/util"
	"github.com/lmittmann/tint"
)

func main() {
	slog.SetDefault(slog.New(
		tint.NewHandler(os.Stderr, &tint.Options{
			Level:      slog.LevelDebug,
			TimeFormat: time.DateTime,
		}),
	))

	conf := bootstrap.ProvideConfiguration()
	router := util.Must(bootstrap.ProvideRouter())

	app := NewApp(conf, router)
	if err := app.Run(); err != nil {
		slog.Error("Application error", "error", err)
	}
}
