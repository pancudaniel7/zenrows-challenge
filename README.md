# Zenrows Challenge

A Go service for managing device templates and user-specific device profiles. The project follows a clean (hexagonal) architecture with clear boundaries between HTTP transport, application use cases, and persistence. Testing spans unit, integration, and acceptance levels using Testcontainers to keep local and CI environments aligned.

## Table of Contents

- [Database](#database)
- [Architecture](#architecture)
- [Prerequisites](#prerequisites)
- [Configuration](#configuration)
- [Local Development](#local-development)
- [Docker & Infrastructure](#docker--infrastructure)
- [Testing](#testing)
- [Project Layout](#project-layout)
- [Error Handling](#error-handling)
- [Authentication](#authentication)
- [Make Targets](#make-targets)
- [Troubleshooting](#troubleshooting)

---

## Database

PostgreSQL is a strong fit for this service because it combines strict relational consistency with the flexibility we need: native uuid identifiers for compact, index-friendly keys, jsonb for evolving device/profile metadata, and robust relational constraints (FKs, uniques) to enforce template–profile relationships.

- **Native UUIDs** – First-class `uuid` type keeps IDs compact and index-friendly (no string hacks).
- **Flexible `JSONB`** – Store evolving device/profile metadata without schema churn while keeping it **queryable** and **indexable**.
- **ACID transactions** – Reliable multi-step updates (e.g., template ↔ profile writes) with serializable isolation when needed.
- **Rich constraints** – Foreign keys, unique and check constraints enforce invariants at the database layer, not just in code.
- **Mature indexing** – B-tree, **GIN/GiST**, partial, and expression indexes cover both relational lookups and JSONB access paths.
- **Extensions & tooling** – `pg_stat_statements`, `pgcrypto`, and battle-tested backup/restore tools streamline observability and ops.

## Architecture

- **Language & Runtime:** Go (modules)  
- **Database:** PostgreSQL (UUID, JSONB, strong transactional model)  
- **Style:** Hexagonal / Ports & Adapters  
- **HTTP Framework:** (via handlers in `internal/adapter/http`)  
- **Persistence:** Repository adapters in `internal/adapter/repo`  
- **Domain Logic:** Use cases in `internal/core/usecase`  
- **Testing:** Unit + integration + acceptance with Testcontainers

This structure keeps dependencies flowing inward. Adapters handle I/O concerns; core use cases encapsulate business logic. Replacing or adding adapters (e.g., alternate storage or transport) is straightforward.

---

## Prerequisites

- Go **1.25**
- Docker and Docker Compose (for Postgres and Testcontainers)
- `make` (optional, for convenience targets)

Verify versions:

```sh
go version
docker --version
docker compose version
make --version
```

---

## Configuration

Configuration is provided via YAML files under `configs/` and can be overridden with environment variables.

- Example profiles:
  - `configs/local.yml`
  - `configs/test.yml`
  - Add more (e.g., `configs/prod.yml`) as needed.

Select a profile by setting `ZENROWS_CONFIG_NAME`:

```sh
export ZENROWS_CONFIG_NAME=local
```

Any setting can be overridden with `ZENROWS_*` environment variables, for example:

```sh
export ZENROWS_DATABASE_HOST=localhost
export ZENROWS_DATABASE_PORT=5432
```

### Example: start with a custom profile

```sh
export ZENROWS_CONFIG_NAME=prod
go run ./cmd
```

> Go modules are resolved automatically. If needed, tidy explicitly:

```sh
go mod tidy
```

---

## Local Development

### 1) Start PostgreSQL (local)

```sh
docker compose -f deployments/docker-compose.yml up -d postgres
```

### 2) Run the API

```sh
go run ./cmd
# or
make run
```

The server listens on the port defined in configuration (e.g., `server.port`).

### 3) Health Check

```sh
curl http://localhost:8080/health
# Expected: "UP!"
```

---

## Docker & Infrastructure

Build the application image:

```sh
make docker-build
# Produces: zenrows-challenge:latest (via build/Dockerfile)
```

Bring up local dependencies (Postgres, etc.):

```sh
make docker-up
```

Tear down:

```sh
make docker-down
```

---

## Testing

### Unit Tests (pure Go, no external deps)

```sh
go test ./internal/...
# or
make unit-tests
```

### Integration & Acceptance Tests

Spins up PostgreSQL 16 via Testcontainers and exercises the real HTTP surface (Fiber handlers).

```sh
GOCACHE=$(pwd)/.gocache go test ./test
# or
make integration-tests
```

### Full Suite

```sh
make tests
```

> **Note:** Docker must be running for integration/acceptance tests.

---

## Project Layout

```
.
├── build/                     # Dockerfiles, build scripts
├── cmd/                       # Application entrypoint (main)
├── configs/                   # YAML configs (local.yml, test.yml, ...)
├── deployments/               # docker-compose, infra manifests
├── internal/
│   ├── adapter/
│   │   ├── http/              # HTTP handlers (transport)
│   │   └── repo/              # Repository adapters (persistence)
│   ├── core/
│   │   └── usecase/           # Application business rules
│   └── pkg/
│       ├── apperr/            # Typed errors (codes/messages/causes)
│       └── middleware/        # Auth and cross-cutting concerns
└── test/                      # Integration & acceptance tests
```

---

## Error Handling

The package `internal/pkg/apperr` defines typed errors with codes, messages, and wrapped causes. Use cases emit these errors; adapters map them to consistent HTTP responses and logs. This yields deterministic client behavior and clear observability.

---

## Authentication

Basic authentication is enforced via middleware (`internal/pkg/middleware`). The middleware:

1. Validates credentials via the authentication service.
2. Injects the caller’s UUID into the request context.
3. Keeps handlers focused on business logic by centralizing credential checks.

---

## Make Targets

Common targets (see `Makefile`):

```sh
make run                # Run the API locally
make unit-tests         # Run unit tests
make integration-tests  # Run integration & acceptance tests (Testcontainers)
make tests              # Run all tests
make docker-build       # Build Docker image (zenrows-challenge:latest)
make docker-up          # Start infra from deployments/docker-compose.yml
```

---

## Troubleshooting

- **Testcontainers fails to start Postgres**
  - Ensure Docker Desktop/Engine is running.
  - Check for port conflicts (e.g., Postgres default port).

- **Health check fails**
  - Confirm `server.port` in your active config.
  - Inspect logs for startup errors.

- **Modules not resolving**
  - Run `go mod tidy`.
  - Remove `$GOMODCACHE` entries if needed and rebuild.

---

## Quick Start

```sh
# 1) Start Postgres
docker compose -f deployments/docker-compose.yml up -d postgres

# 2) Run the app
export ZENROWS_CONFIG_NAME=local
go run ./cmd

# 3) Verify
curl http://localhost:8080/health
```
