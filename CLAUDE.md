# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Go-AI is a restaurant management API backend built with Go using Domain-Driven Design (DDD) patterns. It provides REST APIs for managing 
, menus, media uploads, and user authentication with RBAC.

**Tech Stack:** Go 1.25, Echo v5, PostgreSQL 18, Redis, MinIO, sqlc, JWT

## Build & Development Commands

```bash
# Development
make run              # Run tests then start server with air (hot reload)
make build            # Generate Swagger docs then run server

# Database
make sqlc             # Generate Go code from SQL queries
make up               # Apply database migrations
make down             # Rollback last migration
make migrate          # Create new migration file

# Documentation
make swagger          # Generate Swagger/OpenAPI docs

# Testing
go test ./...         # Run all tests
go test ./internal/container/...  # Run specific package tests

# Full stack
docker-compose up     # Start all services (postgres, redis, minio, app)
```

## Architecture

The codebase follows DDD with Clean Architecture. Each module has four layers:

```
internal/{module}/
├── domain/           # Entities, repository interfaces, value objects
├── application/      # Use cases orchestrating business logic
├── infrastructure/   # Repository implementations, sqlc generated code
│   ├── db/          # Database repository implementations
│   └── sqlc/        # Generated SQL code (do not edit)
└── transport/http/   # HTTP handlers, request/response DTOs
```

### Core Modules

- **identity/** - User auth, JWT tokens, RBAC, sessions (Redis-cached)
- **restaurant/** - Restaurant CRUD with operating hours
- **menu/** - Menus, menu items, option groups, topics
- **media/** - File uploads to MinIO (S3-compatible)
- **health/** - Liveness/readiness probes

### Key Architectural Components

- **Entry point:** `cmd/api/main.go`
- **Server setup:** `internal/app/server.go` (Echo instance, middleware, graceful shutdown)
- **Route registration:** `internal/app/app.go` (all routes under `/api` prefix)
- **Dependency injection:** `internal/container/` (module initialization functions)
- **Middleware stack:** `internal/transport/middlewares/` (CORS, rate limiting, JWT, logging, metrics)
- **Error handling:** `pkg/domain_err/` (AppError with HTTP status mapping)
- **Logging:** `pkg/logger/` (Zerolog with async rotation)

## Database

PostgreSQL with sqlc for type-safe queries. Schema files in `db/schemas/`, migrations in `db/migrations/`.

**sqlc workflow:**
1. Add/modify SQL in `db/schemas/*.sql` or query files
2. Run `make sqlc` to regenerate Go code
3. Generated code goes to `internal/{module}/infrastructure/sqlc/`

## Configuration

Environment variables (or `.env` file):

```
# Required for production
JWT_SECRET, JWT_REFRESH_SECRET

# Database
POSTGRES_HOST, POSTGRES_PORT, POSTGRES_USER, POSTGRES_PASSWORD, POSTGRES_DB

# Redis
REDIS_HOST, REDIS_PORT, REDIS_PASSWORD

# MinIO
MINIO_END_POINT, MINIO_PORT, MINIO_ACCESS_KEY, MINIO_SECRET_KEY, MINIO_BUCKET

# Server
SERVER_PORT (default 8080), ENVIRONMENT (development|production)
```

Config loading: `internal/platform/config/config.go`

## Adding a New Feature

1. Create module structure in `internal/{module}/` with domain, application, infrastructure, transport layers
2. Define entities in `domain/` with repository interface
3. Implement use cases in `application/`
4. Add SQL queries and run `make sqlc` for database operations
5. Create HTTP handlers in `transport/http/`
6. Add initialization in `internal/container/`
7. Register routes in `internal/app/app.go`

## API Documentation

Swagger UI available at `/swagger/index.html` when server is running.

Run `make swagger` to regenerate docs from code annotations.

**Swagger files:** `docs/` (docs.go, swagger.json, swagger.yaml)
**Swagger handler:** `internal/transport/swagger/swagger.go`
