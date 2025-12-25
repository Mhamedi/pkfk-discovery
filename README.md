# PK-FK Discovery: Adapter Studio + Registry + Engine Admin

Enterprise-grade web application for managing database adapters, running PK/FK discovery scans, and administering the adapter registry.

## Architecture

This is a monorepo containing:

- **Backend**: Go 1.25.x REST API (`apps/api`) + Worker (`apps/worker`) using clean architecture
- **Frontend**: Next.js 16.1 with App Router, TypeScript, Tailwind CSS v4, shadcn/ui components
- **Infrastructure**: Postgres (primary DB), Redis (job queue + cache), MinIO (object store)
- **Tooling**: Cross-platform port allocator, Docker Compose with intelligent port selection, CI/CD with Renovate

## Prerequisites

- **Go**: 1.25.x (latest patch) - [Install Go](https://go.dev/doc/install)
- **Node.js**: v24 Active LTS - [Install Node.js](https://nodejs.org/)
- **pnpm**: 10.26.2+ (managed via Corepack)
- **Docker**: Latest stable - [Install Docker](https://docs.docker.com/get-docker/)
- **Docker Compose**: v2.0+ (included with Docker Desktop)

## Quick Start

### 1. Setup Development Environment

```bash
# Install dependencies
make setup

# This will:
# - Check Node.js version (via .nvmrc)
# - Enable Corepack for pnpm
# - Install pnpm dependencies
# - Download Go modules
```

### 2. Start Local Development Stack

```bash
# Start all services with intelligent port allocation
make docker-up

# This will:
# - Run port allocator to find free ports
# - Generate deploy/.env.generated and deploy/docker-compose.override.yml
# - Start all Docker services (web, api, worker, postgres, redis, minio)
# - Run database migrations
# - Seed admin user and demo data
# - Print access URLs and credentials
```

### 3. Access the Application

After `make docker-up` completes, you'll see output like:

```
✓ Services started successfully!

Access URLs:
- Web UI: http://localhost:3000
- API: http://localhost:8080
- MinIO Console: http://localhost:9001

Default Credentials:
- Email: admin@example.com
- Password: admin123
```

## Intelligent Port Allocation

The port allocator (`tools/port-allocator`) automatically detects free ports on your system:

- **Base ports**: WEB=3000, API=8080, PG=5432, REDIS=6379, MINIO=9000, MINIO_CONSOLE=9001
- **Strategy**: If a port is in use, finds the next available port within +100 range
- **Output**: Generates `deploy/.env.generated` and `deploy/docker-compose.override.yml` atomically
- **No manual editing required**: Ports are automatically configured

## Development Commands

```bash
# Format code
make fmt

# Lint code
make lint

# Run tests
make test

# Build binaries and frontend
make build

# View logs
make docker-logs

# Stop services
make docker-down

# Reset everything (removes volumes, generated files)
make docker-reset
```

## Project Structure

```
/
├── apps/
│   ├── api/          # Go REST API
│   ├── worker/       # Go worker/runner
│   └── web/          # Next.js frontend
├── packages/
│   └── shared/       # Shared types/schemas
├── deploy/
│   ├── docker-compose.yml
│   ├── migrations/   # Database migrations
│   └── .env.generated (generated)
├── tools/
│   └── port-allocator/  # Port allocation CLI
└── docs/
    ├── architecture.md
    ├── adapter-spec.md
    └── api.md
```

## Features

- **Adapter Registry**: Catalog and manage database adapters
- **Adapter Studio**: Create, edit, probe, validate, and publish adapters
- **Engine**: Run PK/FK discovery scans with configurable policies
- **Connections**: Manage database connections with encrypted secrets
- **AI Assist**: AI-powered adapter optimization (Ollama/OpenAI-compatible)
- **Administration**: RBAC, audit logs, user management
- **Settings**: Comprehensive system configuration

## Documentation

- [Architecture](./docs/architecture.md) - System design and component diagrams
- [Adapter Specification](./docs/adapter-spec.md) - Adapter bundle structure and requirements
- [API Documentation](./docs/api.md) - REST API reference

## License

Proprietary - All rights reserved

