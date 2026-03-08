# Path to goose binary
GOOSE=$(shell which goose)

# Database connection string
DATABASE_URL?=postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable

# Migration directory
MIGRATION_DIR=internal/adapters/postgres/migrations

generate-oas:
	go run github.com/ogen-go/ogen/cmd/ogen@latest \
		--target internal/adapters/httpapi/gen \
		--clean internal/adapters/httpapi/openapi.yaml

generate-sqlc:
	go run github.com/sqlc-dev/sqlc/cmd/sqlc@latest generate

generate-mappers:
	go run github.com/jmattheis/goverter/cmd/goverter@latest gen github.com/faber-numeris/luciole-auth/authn/internal/platform/mapper

migrate-all-up:
	$(GOOSE) -dir $(MIGRATION_DIR) postgres "$(DATABASE_URL)" up

migrate-all-down:
	$(GOOSE) -dir $(MIGRATION_DIR) postgres "$(DATABASE_URL)" down

.PHONY: generate-oas generate-sqlc generate-mappers migrate-all-up migrate-all-down
