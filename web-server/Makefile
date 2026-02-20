.PHONY: build run dev test test-unit test-integration test-clean lint migrate-up migrate-down migrate-down-all docker-build docker-up docker-down

build:
	go build -o bin/nkrypt-server ./cmd/server

run: build
	@echo "Loading environment from .env file..."
	@set -a && [ -f .env ] && . ./.env && set +a && ./bin/nkrypt-server

dev:
	@echo "Starting development server with live reload (requires air)..."
	@-pkill -f "tmp/main" 2>/dev/null
	@-lsof -ti:9041 | xargs -r kill -9 2>/dev/null
	@set -a && [ -f .env ] && . ./.env && set +a && air

test:
	go test ./... -v

test-unit:
	go clean -testcache
	go test ./internal/... -v

test-integration:
	go clean -testcache
	go test ./test/integration/... -v -tags=integration

lint:
	golangci-lint run

migrate-up:
	@echo "Running migrations (requires golang-migrate CLI)..."
	@if ! command -v migrate >/dev/null 2>&1; then \
		echo "Error: golang-migrate not installed. See README for installation instructions."; \
		exit 1; \
	fi
	@set -a && [ -f .env ] && . ./.env && set +a && migrate -path migrations -database "$$NK_DATABASE_URL" up

migrate-down:
	@echo "Rolling back last migration (requires golang-migrate CLI)..."
	@if ! command -v migrate >/dev/null 2>&1; then \
		echo "Error: golang-migrate not installed. See README for installation instructions."; \
		exit 1; \
	fi
	@set -a && [ -f .env ] && . ./.env && set +a && migrate -path migrations -database "$$NK_DATABASE_URL" down 1

migrate-down-all:
	@echo "Rolling back ALL migrations (requires golang-migrate CLI)..."
	@if ! command -v migrate >/dev/null 2>&1; then \
		echo "Error: golang-migrate not installed. See README for installation instructions."; \
		exit 1; \
	fi
	@set -a && [ -f .env ] && . ./.env && set +a && echo y | migrate -path migrations -database "$$NK_DATABASE_URL" down

docker-build:
	docker build -f deploy/docker/Dockerfile -t nkrypt-server .

docker-up:
	docker compose -f deploy/docker/docker-compose.dev.yml up -d

docker-down:
	docker compose -f deploy/docker/docker-compose.dev.yml down

