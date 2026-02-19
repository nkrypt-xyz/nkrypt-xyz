# Integration Tests

Independent integration tests for the Go web server.

## Prerequisites

- Node.js 20+
- Go server running on `localhost:9041`
- Dependencies: PostgreSQL, Redis, MinIO (started via `make docker-up` in web-server dir)

## Run Tests

```bash
# Install dependencies (first time only)
npm install

# Start server (in another terminal)
cd ../web-server
make docker-up && make migrate-up && make run

# Run tests
./test-go-compliance.sh
```
