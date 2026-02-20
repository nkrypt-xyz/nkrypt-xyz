# nkrypt.xyz Web Server

High-performance Go web server for encrypted file storage.

## Quick Start

### 1. Prerequisites

- **Go 1.21+** - [Install Go](https://go.dev/doc/install)
- **Docker/Podman** - For dependencies (PostgreSQL, Redis, MinIO)
- **Make** - For build automation
- **golang-migrate** (optional) - For database migrations
- **air** (optional) - For live reload during development

### 2. Start Dependencies

```bash
make docker-up
```

Starts PostgreSQL (5432), Redis (6379), and MinIO (9000/9001).

### 3. Run Migrations

```bash
# Option 1: Using golang-migrate (recommended)
make migrate-up

> **Note**: `make migrate-up` requires `golang-migrate` CLI. 

# Option 2: Manual SQL piping (no migrate CLI needed)
for f in migrations/*.up.sql; do 
  docker exec -i nkrypt-postgres-dev psql -U nkrypt -d nkrypt < "$f"
done
```

### 4. Configure Environment

Copy `.env.example` to `.env`:

```bash
cp .env.example .env
```

The `.env.example` file is already configured for local development with default values. 

**Required variables:**
- `NK_DATABASE_URL` - PostgreSQL connection string
- `NK_REDIS_ADDR` - Redis address
- `NK_MINIO_ENDPOINT`, `NK_MINIO_ACCESS_KEY`, `NK_MINIO_SECRET_KEY` - MinIO config
- `NK_IAM_DEFAULT_ADMIN_PASSWORD` - Admin password

> **Note**: All `make` commands automatically load variables from `.env`

See `.env.example` for all available configuration options.

### 5. Run Server

```bash
make run
```

Server starts on `http://localhost:9041`

### 6. Test

Login as admin:

```bash
curl -X POST http://localhost:9041/api/user/login \
  -H "Content-Type: application/json" \
  -d '{"userName":"admin","password":"PleaseChangeMe@YourEarliest2Day"}'
```

Check health:

```bash
curl http://localhost:9041/healthz
```

Run integration tests:

```bash
cd ../independent-integration-test
npm install
npm run test:compliance
```

## Development

### Make Commands

```bash
make build              # Build binary
make run                # Build and run (loads .env)
make dev                # Run with live reload (requires air)
make test-unit          # Run unit tests
make test-integration   # Run integration tests (needs running server)
make docker-up          # Start dependencies
make docker-down        # Stop dependencies
make migrate-up         # Run migrations (requires golang-migrate)
make migrate-down       # Rollback migration (requires golang-migrate)
```

### Live Reload (Development)

For automatic reload on code changes:

```bash
# Install air (one-time)
go install github.com/cosmtrek/air@latest

# Run with live reload
make dev
```

This watches for file changes and automatically rebuilds/restarts the server.

### Project Structure

```
cmd/server/          # Application entrypoint
internal/
  ├── config/        # Configuration (Viper)
  ├── handler/       # HTTP handlers
  ├── middleware/    # Auth, logging, CORS
  ├── model/         # Domain models
  ├── pkg/           # Shared utilities
  ├── repository/    # Data access (SQL)
  ├── router/        # Route definitions
  ├── server/        # HTTP server
  └── service/       # Business logic
migrations/          # Database migrations
test/                # Integration tests
```

## Documentation

- **[API Reference](../dev-docs/API.md)** - Complete endpoint documentation
- **[Architecture](../dev-docs/ARCHITECTURE.md)** - System design and patterns
- **[Database Schema](../dev-docs/DATABASE.md)** - Tables, migrations, and relationships
- **[Contributing](../dev-docs/CONTRIBUTING.md)** - Development guidelines


---

**Need help?** See [dev-docs/](../dev-docs/) for detailed guides.
