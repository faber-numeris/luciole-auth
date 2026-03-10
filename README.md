# luciole-auth

AuthN and AuthZ components used on the luciole ecosystem.

## Project Overview

`luciole-auth` is the central authentication and authorization service for the Luciole ecosystem. It currently focuses on authentication (AuthN) through the `authn` service.

The project follows a **Hexagonal Architecture** (Ports and Adapters) to ensure a clean separation between business logic and infrastructure concerns.

## Hexagonal Architecture

This project is organized into layers to maintain a strict dependency rule: **inner layers never depend on outer layers.**

1.  **Domain (`internal/domain`)**: Pure business logic and entities. Zero dependencies on other internal packages or external frameworks.
2.  **Application (`internal/app`)**: Orchestrates business use cases. It depends only on the Domain and defines **Ports** (interfaces) for the infrastructure it needs (e.g., repositories, mailers).
3.  **Ports (`internal/app/ports`)**: Interfaces defined by the application layer that must be implemented by the infrastructure/adapters layer.
4.  **Adapters (`internal/adapters`)**: Implementations of the Ports.
    *   **httpapi**: Implements the HTTP/OpenAPI server (using `ogen`).
    *   **postgres**: Implements the repository interfaces using `sqlc` and PostgreSQL.
    *   **mail**: Implements the mailer interface.
5.  **Infrastructure (`internal/infrastructure`)**: General-purpose technical concerns like configuration loading and database connection management.
6.  **Cmd (`cmd/authn`)**: The entry point. `app.go` in this directory is the **only** place where all layers are wired together (Dependency Injection).

## Technical Stack

- **Language:** [Go](https://go.dev/) (1.24+)
- **Database:** [PostgreSQL](https://www.postgresql.org/)
- **API Specification:** [OpenAPI 3.0](https://swagger.io/specification/)
- **Code Generation:**
  - [sqlc](https://sqlc.dev/) for type-safe SQL queries.
  - [ogen](https://ogen.dev/) for OpenAPI-based code generation.
  - [goverter](https://github.com/jmattheis/goverter) for type-safe struct conversion.
- **Migrations:** [goose](https://github.com/pressly/goose) for database versioning.

## Project Structure

```
authn/
├── cmd/authn/              # Entry point & dependency wiring
├── internal/
│   ├── domain/             # Pure entities (User, etc.)
│   ├── app/                # Business logic & services
│   │   └── ports/          # Interfaces (Repository, Mailer)
│   ├── adapters/           # External implementations
│   │   ├── httpapi/        # HTTP Handlers & OpenAPI
│   │   ├── postgres/       # SQLC Repositories & Migrations
│   │   └── mail/           # SMTP/Mailpit implementations
│   ├── infrastructure/     # Config & DB connection init
│   └── platform/           # Shared utilities (Mappers, etc.)
├── Makefile                # Generation and migration tasks
└── sqlc.yaml               # sqlc configuration
```

## Development Workflow

### Code Generation

We use several generators to maintain type safety:

```bash
# Generate API code from OpenAPI spec
make generate-oas

# Generate type-safe SQL code from queries
make generate-sqlc

# Generate type-safe struct mappers
make generate-mappers
```

### Database Migrations

Migrations are located in `authn/internal/adapters/outbound/postgres/migrations`.

```bash
# Run all migrations up
make migrate-all-up

# Revert all migrations
make migrate-all-down
```

### Running the Application

```bash
cd authn
go run ./cmd/authn
```
