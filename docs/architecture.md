# Architecture Documentation

## System Overview

PK-FK Discovery is an enterprise-grade monorepo application for managing database adapters, running PK/FK discovery scans, and administering the adapter registry.

## Architecture Diagram

```mermaid
graph TB
    subgraph Frontend["Frontend Layer"]
        Web[Next.js Web App]
    end
    
    subgraph Backend["Backend Layer"]
        API[Go REST API]
        Worker[Go Worker]
    end
    
    subgraph Infrastructure["Infrastructure Layer"]
        Postgres[(PostgreSQL)]
        Redis[(Redis)]
        MinIO[(MinIO)]
    end
    
    Web -->|HTTP/REST| API
    API -->|Job Queue| Redis
    Worker -->|Job Queue| Redis
    API -->|Persistence| Postgres
    Worker -->|Persistence| Postgres
    API -->|Object Storage| MinIO
    Worker -->|Object Storage| MinIO
    
    Worker -->|SQL Queries| ExternalDB[(External Databases)]
```

## Component Architecture

### Backend (Go)

#### Clean Architecture Layers

```
apps/api/
├── cmd/api/main.go              # Application entry point
├── internal/
│   ├── transport/http/          # HTTP handlers, routing, middleware
│   ├── domain/                  # Business entities and interfaces
│   ├── services/                # Business logic and use cases
│   ├── repositories/            # Data access layer (Postgres)
│   └── integrations/            # External services (Redis, MinIO)
│       ├── redis/
│       └── minio/
```

#### Key Components

- **Transport Layer**: Chi router, REST handlers, middleware (auth, CORS, logging)
- **Domain Layer**: Core entities (User, Adapter, Connection, Scan, Job, AuditLog)
- **Service Layer**: Business logic orchestration
- **Repository Layer**: Postgres implementations
- **Integration Layer**: Redis (job queue, caching), MinIO (object storage)

### Frontend (Next.js)

#### Structure

```
apps/web/
├── app/                         # Next.js App Router pages
├── components/                  # React components
│   ├── ui/                      # shadcn/ui components
│   └── layout/                  # Layout components
├── lib/                         # Utilities and API client
└── hooks/                       # React hooks
```

## Data Flow

### Adapter Creation Flow

```mermaid
sequenceDiagram
    participant User
    participant Web
    participant API
    participant Worker
    participant Redis
    participant Postgres
    participant MinIO
    
    User->>Web: Create Draft
    Web->>API: POST /studio/drafts
    API->>Postgres: Save Draft
    API->>Web: Return Draft ID
    
    User->>Web: Probe Connection
    Web->>API: POST /studio/drafts/{id}/probe
    API->>Redis: Enqueue Probe Job
    API->>Web: Return Job ID
    
    Worker->>Redis: Poll Jobs
    Worker->>ExternalDB: Execute Probe SQL
    Worker->>Postgres: Update Draft with Results
    
    User->>Web: Validate Adapter
    Web->>API: POST /studio/drafts/{id}/validate
    API->>Redis: Enqueue Validation Job
    Worker->>Postgres: Run Validation Tests
    Worker->>Postgres: Update Results
    
    User->>Web: Publish Adapter
    Web->>API: POST /studio/drafts/{id}/publish
    API->>API: Package Bundle
    API->>API: Generate Signature
    API->>MinIO: Upload Bundle
    API->>Postgres: Register Adapter Version
```

### Scan Execution Flow

```mermaid
sequenceDiagram
    participant User
    participant Web
    participant API
    participant Worker
    participant Redis
    participant Postgres
    participant ExternalDB
    
    User->>Web: Create Scan
    Web->>API: POST /scans
    API->>Postgres: Create Scan Record
    API->>Redis: Enqueue Scan Job
    API->>Web: Return Scan ID
    
    Worker->>Redis: Poll Jobs
    Worker->>Postgres: Load Adapter & Connection
    Worker->>ExternalDB: Metadata Extraction
    Worker->>ExternalDB: Profiling
    Worker->>ExternalDB: Evidence Collection
    Worker->>Worker: Candidate Generation
    Worker->>Worker: Scoring & Graph Reconciliation
    Worker->>Postgres: Save Results
    Worker->>Redis: Update Job Status
```

## Security Architecture

### Authentication & Authorization

- **JWT-based authentication** with refresh tokens
- **RBAC**: Admin, Editor, Viewer roles enforced server-side
- **Future-proof**: Interfaces for SSO/LDAP integration

### Data Protection

- **Secrets encryption**: AES-256-GCM at rest (connections, AI API keys)
- **SQL safety**: Worker enforces SELECT/EXPLAIN only, keyword denylist
- **Audit logging**: All sensitive actions logged with user, IP, timestamp

## Deployment Architecture

### Docker Compose Services

- **web**: Next.js frontend (port 3000)
- **api**: Go REST API (port 8080)
- **worker**: Go worker (no exposed port)
- **postgres**: PostgreSQL (port 5432)
- **redis**: Redis (port 6379)
- **minio**: MinIO object storage (ports 9000, 9001)

### Port Allocation

Intelligent port allocator (`tools/port-allocator`) detects free ports and generates override files automatically.

## Observability

### Logging

- **Structured JSON logs** (logrus/zap)
- **Request IDs** for tracing
- **Audit logs** for compliance

### Health Checks

- `/healthz`: Liveness probe
- `/readyz`: Readiness probe (checks DB/Redis/MinIO)
- `/metrics`: Prometheus metrics

## Scalability Considerations

- **Horizontal scaling**: Stateless API, multiple workers
- **Job queue**: Redis-based with concurrency limits
- **Caching**: Redis for frequently accessed data
- **Object storage**: MinIO for large artifacts

