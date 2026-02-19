package di

import (
	"fmt"
	"log/slog"
	"net/http"

	api "github.com/faber-numeris/luciole-auth/api/gen"
	"github.com/faber-numeris/luciole-auth/authn/configuration"
	"github.com/faber-numeris/luciole-auth/authn/handlers"
	"github.com/faber-numeris/luciole-auth/authn/persistence/database"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	specui "github.com/oaswrap/spec-ui"
)

func ProvideConfiguration() configuration.IAppConfig {
	cfg, err := configuration.Load()
	if err != nil {
		panic(err)
	}
	return cfg
}

func ProvideAuthnService() handlers.IAuthnService {
	srv := handlers.NewAuthnService()
	return srv
}

func ProvideSecurityService() handlers.ISecurityService {
	srv := handlers.NewSecurityService()
	return srv
}
func ProvideHandler() (http.Handler, error) {
	securityService := ProvideSecurityService()
	srv := ProvideAuthnService()
	return api.NewServer(srv, securityService)
}

func ProvideRouter() (http.Handler, error) {
	// -------------------------------
	// Redoc documentation UI
	// -------------------------------
	// Stoplight Elements
	specuiHandler := specui.NewHandler(
		specui.WithTitle("Luciole Auth API"),
		specui.WithDocsPath("/docs/authn"),
		specui.WithSpecPath("/docs/authn/openapi.yaml"),
		specui.WithSpecFile("api/openapi.yaml"),
		specui.WithStoplightElements(),
	)

	mux := chi.NewRouter()

	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)
	mux.Get(specuiHandler.DocsPath(), specuiHandler.DocsFunc())
	mux.Get(specuiHandler.SpecPath(), specuiHandler.SpecFunc())

	srv, err := ProvideHandler()
	if err != nil {
		return nil, err
	}

	mux.Mount("/v1/", http.StripPrefix("/v1", srv))
	err = chi.Walk(
		mux,
		func(
			method string,
			route string,
			handler http.Handler,
			middlewares ...func(http.Handler) http.Handler,
		) error {
			slog.Info("Registering route", slog.String("method", method), slog.String("route", route))
			return nil
		},
	)

	if err != nil {
		return nil, fmt.Errorf("failed to walk routes: %w", err)
	}
	return mux, nil
}

func ProvideDatabase() (*database.DB, error) {
	cfg := ProvideConfiguration()
	return database.NewConnection(cfg)
}
