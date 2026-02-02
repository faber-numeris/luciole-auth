# Path to goose binary (if not already on PATH)
GOOSE=$(shell which goose)

# Database connection string
DB_CONNECTION=postgres://postgres:postgres@localhost/postgres?sslmode=disable

# Migration directory
MIGRATION_DIR=migration

# Ensure required args are provided
require-name:
	@if [ -z "$(name)" ]; then \
		echo "❌ Missing migration name. Usage:"; \
		echo "   make create-migration name=init_users"; \
		exit 1; \
	fi

# Create a new migration file
create-migration: require-name
	$(GOOSE) -dir $(MIGRATION_DIR) create $(name) sql

# Run migrations by one step (up)
migrate-one-up:
	$(GOOSE) -dir $(MIGRATION_DIR) postgres "$(DB_CONNECTION)" up 1

# Run all migrations up to the latest
migrate-all-up:
	$(GOOSE) -dir $(MIGRATION_DIR) postgres "$(DB_CONNECTION)" up

# Revert one migration step (down)
migrate-one-down:
	$(GOOSE) -dir $(MIGRATION_DIR) postgres "$(DB_CONNECTION)" down 1

# Revert all migrations to the initial state
migrate-all-down:
	$(GOOSE) -dir $(MIGRATION_DIR) postgres "$(DB_CONNECTION)" down

clean-sqlc:
	rm -rf persistence/sqlc/*.go

generate-sqlc: clean-sqlc
	go run github.com/sqlc-dev/sqlc/cmd/sqlc@latest generate


.PHONY: create-migration migrate-one-up migrate-all-up migrate-one-down migrate-all-down clean-migrations generate-sqlc clean-sqlc

