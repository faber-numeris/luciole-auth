package bootstrap

import (
	"fmt"

	postgresadapter "github.com/faber-numeris/luciole-auth/authn/internal/adapters/outbound/postgres"
	"github.com/faber-numeris/luciole-auth/authn/internal/adapters/outbound/postgres/gen"
	inboundport "github.com/faber-numeris/luciole-auth/authn/internal/app/ports/inbound"
	outboundport "github.com/faber-numeris/luciole-auth/authn/internal/app/ports/outbound"
	"github.com/faber-numeris/luciole-auth/authn/internal/app/service"
	"github.com/faber-numeris/luciole-auth/authn/internal/infrastructure/config"
	infrapostgres "github.com/faber-numeris/luciole-auth/authn/internal/infrastructure/postgres"
)

func ProvideHashingService() inboundport.HashingService {
	return service.NewHashingService()
}

func ProvideRepository() outboundport.Repository {

	cfg, err := config.LoadConfig()
	if err != nil {
		panic(fmt.Errorf("failed to load configuration: %w", err))
	}

	pool := infrapostgres.Connect(cfg)
	var repo = &struct {
		outboundport.UserRepository
		outboundport.UserConfirmationRepository
	}{
		postgresadapter.NewUserRepository(gen.New(pool)),
		postgresadapter.NewUserConfirmationRepository(gen.New(pool)),
	}

	return repo
}
