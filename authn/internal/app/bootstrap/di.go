package bootstrap

import (
	"fmt"

	postgresadapter "github.com/faber-numeris/luciole-auth/authn/internal/adapters/outbound/postgres"
	inboundport "github.com/faber-numeris/luciole-auth/authn/internal/ports/inbound"
	outboundport "github.com/faber-numeris/luciole-auth/authn/internal/ports/outbound"
	"github.com/faber-numeris/luciole-auth/authn/internal/core/services"
	"github.com/faber-numeris/luciole-auth/authn/internal/infrastructure/config"
	infrapostgres "github.com/faber-numeris/luciole-auth/authn/internal/infrastructure/postgres"
)

func ProvideHashingService() inboundport.HashingService {
	return services.NewHashingService()
}

func ProvideRepository() outboundport.Repository {

	_, err := config.LoadConfig()
	if err != nil {
		panic(fmt.Errorf("failed to load configuration: %w", err))
	}

	db := infrapostgres.Connect()
	var repo = &struct {
		outboundport.UserRepository
		outboundport.UserConfirmationRepository
	}{
		postgresadapter.NewUserRepository(db),
		postgresadapter.NewUserConfirmationRepository(db),
	}

	return repo
}
