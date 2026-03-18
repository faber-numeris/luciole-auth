package bootstrap

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

	"github.com/faber-numeris/luciole-auth/authn/internal/adapters/inbound/httpapi"
	api "github.com/faber-numeris/luciole-auth/authn/internal/adapters/inbound/httpapi/gen"
	"github.com/faber-numeris/luciole-auth/authn/internal/adapters/outbound/mail"
	"github.com/faber-numeris/luciole-auth/authn/internal/core/services"
	"github.com/faber-numeris/luciole-auth/authn/internal/infrastructure/config"
	"github.com/faber-numeris/luciole-auth/authn/internal/infrastructure/postgres"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jmoiron/sqlx"
	specui "github.com/oaswrap/spec-ui"
)

type App struct {
	server *http.Server
	db     *sqlx.DB
}

func NewApp() *App {
	mailer := mail.NewService()
	hashingSvc := ProvideHashingService()
	repo := ProvideRepository()
	userSvc := services.NewUserService(repo, hashingSvc, mailer)
	encryptionSvc := services.NewEncryptionService()

	handler := httpapi.NewHandler(userSvc, hashingSvc, encryptionSvc)
	security := httpapi.NewSecurityHandler()

	srv, err := api.NewServer(handler, security)
	if err != nil {
		panic(fmt.Errorf("failed to create server: %w", err))
	}

	router := buildRouter(srv)

	cfg, err := config.LoadConfig()
	if err != nil {
		panic(fmt.Errorf("failed to load configuration: %w", err))
	}
	address := fmt.Sprintf(":%d", cfg.Port())
	return &App{
		server: &http.Server{
			Addr:    address,
			Handler: router,
		},
		db: postgres.Connect(),
	}
}

func buildRouter(srv *api.Server) http.Handler {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(fmt.Errorf("failed to load configuration: %w", err))
	}
	specuiHandler := specui.NewHandler(
		specui.WithTitle("Luciole Auth API"),
		specui.WithDocsPath("/docs/authn"),
		specui.WithSpecPath("/docs/authn/openapi.yaml"),
		specui.WithSpecFile("authn/internal/adapters/inbound/httpapi/openapi.yaml"),
		specui.WithStoplightElements(),
	)

	mux := chi.NewRouter()

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   cfg.AllowedOrigins(),
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Requested-With"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)
	mux.Get(specuiHandler.DocsPath(), specuiHandler.DocsFunc())
	mux.Get(specuiHandler.SpecPath(), specuiHandler.SpecFunc())

	mux.Mount("/v1/", http.StripPrefix("/v1", srv))

	return mux
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

		if err := a.db.Close(); err != nil {
			slog.Error("failed to close database", "error", err)
		}

		slog.Info("Server exited")
	}

	return nil
}
