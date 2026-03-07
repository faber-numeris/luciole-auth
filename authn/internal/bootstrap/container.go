package bootstrap

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/faber-numeris/luciole-auth/authn/internal/adapters/http/api/gen"
	"github.com/faber-numeris/luciole-auth/authn/internal/adapters/http/handlers"
	"github.com/faber-numeris/luciole-auth/authn/internal/adapters/mail"
	"github.com/faber-numeris/luciole-auth/authn/internal/adapters/persistence/postgres"
	"github.com/faber-numeris/luciole-auth/authn/internal/adapters/persistence/postgres/sqlc"
	"github.com/faber-numeris/luciole-auth/authn/internal/application"
	"github.com/faber-numeris/luciole-auth/authn/internal/bootstrap/config"
	"github.com/faber-numeris/luciole-auth/authn/internal/ports/messaging"
	"github.com/faber-numeris/luciole-auth/authn/internal/ports/repository"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	specui "github.com/oaswrap/spec-ui"
)

func ProvideConfiguration() config.IAppConfig {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}
	return cfg
}

func ProvideUserRepository() repository.IUserRepository {
	// Get Database Singleton
	db, err := config.GetInstance(ProvideConfiguration())
	if err != nil {
		if db != nil {
			_ = db.Close()
		}
		panic(err)
	}

	repo := postgres.NewSQLCUserRepository(sqlc.New(db.Pool))

	return repo
}

func ProvideUserConfirmationRepository() repository.IUserConfirmationRepository {
	// Get Database Singleton
	db, err := config.GetInstance(ProvideConfiguration())
	if err != nil {
		if db != nil {
			_ = db.Close()
		}
		panic(err)
	}

	repo := postgres.NewSQLCUserConfirmationRepository(sqlc.New(db.Pool))

	return repo
}

func ProvideHashingService() application.IHashingService {
	return application.NewHashingService()
}

func ProvideMailService() messaging.IMailService {
	cfg := ProvideConfiguration()
	return mail.NewMailpit(cfg)
}

func ProvideUserService() application.IUserService {
	return application.NewUserService(
		ProvideUserRepository(),
		ProvideUserConfirmationRepository(),
		ProvideHashingService(),
		ProvideMailService(),
	)
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
	specuiHandler := specui.NewHandler(
		specui.WithTitle("Luciole Auth API"),
		specui.WithDocsPath("/docs/authn"),
		specui.WithSpecPath("/docs/authn/openapi.yaml"),
		specui.WithSpecFile("api/openapi.yaml"),
		specui.WithStoplightElements(),
	)

	mux := chi.NewRouter()

	cfg := ProvideConfiguration()
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
