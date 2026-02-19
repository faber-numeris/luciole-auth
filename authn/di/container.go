package di

import (
	"fmt"
	"log/slog"
	"net/http"

	api "github.com/faber-numeris/luciole-auth/api/gen"
	"github.com/faber-numeris/luciole-auth/authn/configuration"
	"github.com/faber-numeris/luciole-auth/authn/handlers"
	"github.com/faber-numeris/luciole-auth/authn/persistence/database"
	"github.com/faber-numeris/luciole-auth/authn/persistence/repository"
	sqlc2 "github.com/faber-numeris/luciole-auth/authn/persistence/sqlc"
	"github.com/faber-numeris/luciole-auth/authn/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	specui "github.com/oaswrap/spec-ui"
)

func ProvideConfiguration() configuration.IAppConfig {
	cfg, err := configuration.Load()
	if err != nil {
		panic(err)
	}
	return cfg
}

func ProvideDatabase() (*database.DB, error) {
	return database.GetInstance(ProvideConfiguration())
}

func ProvideUserRepository() repository.IUserRepository {
	db, err := ProvideDatabase()
	if err != nil {
		if db != nil {
			_ = db.Close()
		}
		panic(err)
	}

	sqlc := sqlc2.New(db.Pool)

	repo := repository.NewSQLCUserRepository(sqlc)

	return repo
}

func ProvideUserService() service.IUserService {

	return service.NewUserService(ProvideUserRepository())
}

func ProvideAuthnService() handlers.IAuthnService {
	userService := ProvideUserService()
	srv := handlers.NewAuthnService(userService)
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

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:*", "https://*"},
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
