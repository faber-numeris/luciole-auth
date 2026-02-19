# luciole-auth

AuthN and AuthZ components used on the luciole ecosystem.

## Project Overview

`luciole-auth` is the central authentication and authorization service for the Luciole ecosystem. It currently focus on authentication (AuthN) through the `authn` service.

## Technical Stack

- **Language:** [Go](https://go.dev/) (1.24.8+)
- **Database:** [PostgreSQL](https://www.postgresql.org/)
- **Infrastructure:** [Docker](https://www.docker.com/) & [Docker Compose](https://docs.docker.com/compose/)
- **API Specification:** [OpenAPI 3.0](https://swagger.io/specification/) (located at `api/openapi.yaml`)
- **Code Generation:**
  - [sqlc](https://sqlc.dev/) for type-safe SQL queries.
  - [ogen](https://ogen.dev/) for OpenAPI-based code generation.
  - [goverter](https://github.com/jmattheis/goverter) for type-safe struct conversion.
- **Migrations:** [goose](https://github.com/pressly/goose) for database versioning.

## Development Setup

### Prerequisites

- **Go** 1.24.8 or later
- **Docker** and **Docker Compose**
- **Make** (for running build scripts)

### Required Tools

Install the following Go-based tools:

```bash
# Install sqlc for type-safe SQL generation
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# Install ogen for OpenAPI code generation  
go install github.com/ogen-go/ogen/cmd/ogen@latest

# Install goose for database migrations
go install github.com/pressly/goose/v3/cmd/goose@latest

# Install goverter for type-safe struct conversion
go install github.com/jmattheis/goverter/cmd/goverter@v1.9.2
```

Make sure `$(go env GOPATH)/bin` is in your `PATH`.

### Database Setup

1. Start the PostgreSQL database using Docker Compose:
   ```bash
   docker compose up -d
   ```
   *Note: This starts PostgreSQL with default credentials (`postgres:postgres` on `localhost:5432`).*

2. Create the database:
   ```bash
   createdb postgres
   ```

### Development Workflow

#### Database Migrations

```bash
# Create a new migration
make create-migration name=migration_name

# Run all migrations
make migrate-all-up

# Run one migration
make migrate-one-up

# Revert one migration
make migrate-one-down

# Revert all migrations
make migrate-all-down
```

#### Code Generation

```bash
# Generate type-safe SQL code from queries
make generate-sqlc

# Generate API code from OpenAPI spec (located at api/openapi.yaml)
make generate
```

#### Running the Application

```bash
cd authn
go run ./cmd/main.go
```

### Environment Configuration

The application uses environment variables for configuration. You can create a `.env` file in the `authn` directory.

## Project Structure

```
luciole-auth/
├── api/
│   ├── gen/                # Generated API code from OpenAPI spec
│   └── openapi.yaml        # OpenAPI specification (API definitions)
├── authn/                  # Main authentication service
│   ├── cmd/                # Application entry point
│   ├── configuration/      # Configuration handling
│   ├── di/                 # Dependency injection
│   ├── handlers/           # API handlers (OpenAPI implementation)
│   ├── migration/          # Database migrations (SQL files)
│   ├── model/              # Domain models and conversions
│   ├── persistence/        # Data access layer
│   │   ├── database/       # DB connection handling
│   │   ├── queries/        # SQL queries for sqlc
│   │   ├── repository/     # Repository pattern implementation
│   │   └── sqlc/           # Generated SQL code
│   ├── service/            # Business logic (services)
│   └── sqlc.yaml           # sqlc configuration
├── docker-compose.yaml     # Infrastructure setup (PostgreSQL)
├── Makefile                # Build and generation tasks
├── go.mod                  # Go module definition
└── README.md               # Project documentation
```

## Roadmap

- [ ] Develop a lightweight **IAM server** to handle Identity and Authentication (AuthN).
- [ ] Integrate and leverage **[OpenFGA](https://openfga.dev/)** as the core Authorization engine (AuthZ) for fine-grained access control.
