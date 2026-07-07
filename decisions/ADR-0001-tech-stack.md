# ADR-0001: Initial Technology Stack

## Status

Accepted

## Context

The project should demonstrate backend engineering skill while remaining buildable by one person with Codex assistance. The previous project used Go, Gin, SQLite, migrations, JWT, bcrypt, and Swagger.

## Decision

Use Go for the backend and PostgreSQL for persistence. Use Docker Compose for local dependencies. Use OpenAPI for API documentation. Use `chi` as the HTTP router.

Preferred stack:

- Go 1.23 or newer.
- PostgreSQL 16 or newer.
- `chi` for HTTP routing.
- `pgx` for database access.
- `golang-migrate` for migrations.
- OpenAPI 3.1 for API docs.
- Docker Compose for local development.
- GitHub Actions for CI.

`chi` is a good fit for the foundation because it keeps handlers close to standard `net/http`, has a small API surface, and provides common middleware without requiring a full web framework. The query strategy is not finalized in this ADR; `sqlc` remains preferred for the database phase, with handwritten `pgx` repositories acceptable only with explicit justification.

## Consequences

Positive:

- Go is strong for backend APIs and portfolio signaling.
- PostgreSQL allows realistic constraints, transactions, and indexing.
- Docker Compose improves reviewer experience.
- OpenAPI makes the API contract visible.
- `chi` keeps routing explicit and easy to test with standard library tools.

Tradeoffs:

- More setup complexity than SQLite.
- More decisions around migrations and integration tests.
- Requires careful scope control.
- `chi` provides fewer framework-level conventions than Gin, so project package boundaries need to stay disciplined.

## Review Notes

Router choice was resolved during Phase 1 bootstrap on 2026-07-07. Database query generation remains a Phase 2 decision, with `sqlc` as the current preference from `project.yaml`.
