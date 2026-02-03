package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/faber-numeris/luciole-auth/authn/di"
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
	err = http.ListenAndServe(address, router)
	if err != nil {
		panic(err)
	}
}
