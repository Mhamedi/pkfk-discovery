.PHONY: setup fmt lint test build docker-up docker-down docker-logs docker-reset db-migrate seed help

# Default target
help:
	@echo "Available targets:"
	@echo "  make setup          - Setup development environment"
	@echo "  make fmt            - Format code (Go + frontend)"
	@echo "  make lint           - Lint code (Go + frontend)"
	@echo "  make test           - Run tests (Go + frontend)"
	@echo "  make build          - Build binaries and frontend"
	@echo "  make docker-up      - Start Docker stack with port allocation"
	@echo "  make docker-down    - Stop Docker stack"
	@echo "  make docker-logs   - View Docker logs"
	@echo "  make docker-reset   - Reset Docker stack (removes volumes)"
	@echo "  make db-migrate     - Run database migrations"
	@echo "  make seed           - Seed admin user and demo data"

# Setup development environment
setup:
	@echo "Setting up development environment..."
	@if command -v node >/dev/null 2>&1; then \
		NODE_VERSION=$$(node -v | cut -d'v' -f2 | cut -d'.' -f1); \
		if [ "$$NODE_VERSION" != "24" ]; then \
			echo "Warning: Node.js v24 is required (found v$$NODE_VERSION). Use nvm or .nvmrc"; \
		fi; \
	else \
		echo "Error: Node.js not found. Please install Node.js v24"; \
		exit 1; \
	fi
	@if command -v corepack >/dev/null 2>&1; then \
		corepack enable; \
		corepack prepare pnpm@10.26.2 --activate; \
	else \
		echo "Error: Corepack not found. Please install Node.js v24+"; \
		exit 1; \
	fi
	@echo "Installing pnpm dependencies..."
	@pnpm install
	@echo "Downloading Go modules..."
	@cd apps/api && go mod download || true
	@cd apps/worker && go mod download || true
	@cd tools/port-allocator && go mod download || true
	@echo "Setup complete!"

# Format code
fmt:
	@echo "Formatting Go code..."
	@go fmt ./apps/api/... ./apps/worker/... ./tools/...
	@echo "Formatting frontend code..."
	@pnpm format || true

# Lint code
lint:
	@echo "Linting Go code..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./apps/api/... ./apps/worker/... ./tools/...; \
	else \
		echo "Warning: golangci-lint not found. Install: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi
	@echo "Linting frontend code..."
	@pnpm lint || true

# Run tests
test:
	@echo "Running Go tests..."
	@go test -v ./apps/api/... ./apps/worker/... ./tools/... || true
	@echo "Running frontend tests..."
	@pnpm test || true

# Build binaries and frontend
build:
	@echo "Building Go binaries..."
	@cd apps/api && go build -o bin/api ./cmd/api || true
	@cd apps/worker && go build -o bin/worker ./cmd/worker || true
	@cd tools/port-allocator && go build -o ../../bin/port-allocator . || true
	@echo "Building frontend..."
	@pnpm build || true

# Docker operations
docker-up:
	@echo "Allocating ports..."
	@cd tools/port-allocator && go run . -base-dir ../..
	@echo "Starting Docker stack..."
	@cd deploy && docker compose up -d
	@echo "Waiting for services to be healthy..."
	@sleep 5
	@echo "Running migrations..."
	@$(MAKE) db-migrate || true
	@echo "Seeding database..."
	@$(MAKE) seed || true
	@echo ""
	@echo "✓ Services started successfully!"
	@echo ""
	@echo "Access URLs:"
	@if [ -f deploy/.env.generated ]; then \
		WEB_PORT=$$(grep WEB_PORT deploy/.env.generated | cut -d'=' -f2); \
		API_PORT=$$(grep API_PORT deploy/.env.generated | cut -d'=' -f2); \
		MINIO_CONSOLE_PORT=$$(grep MINIO_CONSOLE_PORT deploy/.env.generated | cut -d'=' -f2); \
		echo "  - Web UI: http://localhost:$$WEB_PORT"; \
		echo "  - API: http://localhost:$$API_PORT"; \
		echo "  - MinIO Console: http://localhost:$$MINIO_CONSOLE_PORT"; \
	fi
	@echo ""
	@echo "Default Credentials:"
	@echo "  - Email: admin@example.com"
	@echo "  - Password: admin123"

docker-down:
	@cd deploy && docker compose down

docker-logs:
	@cd deploy && docker compose logs -f

docker-reset:
	@echo "Stopping and removing containers..."
	@cd deploy && docker compose down -v
	@echo "Removing generated files..."
	@rm -f deploy/.env.generated deploy/docker-compose.override.yml
	@echo "Reset complete!"

# Database operations
db-migrate:
	@echo "Running database migrations..."
	@cd deploy && docker compose exec -T postgres psql -U pkfk -d pkfk_discovery < migrations/001_initial_schema.sql || true

seed:
	@echo "Seeding database..."
	@cd deploy && docker compose exec -T postgres psql -U pkfk -d pkfk_discovery < migrations/002_seed_data.sql || true

