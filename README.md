# GatherOps

GatherOps is a production-style event operations API for organizations that need to create events, manage attendees, collect RSVPs, and audit operational changes. This repository contains the planning, specification, prompt workflow, and implementation workspace for building it as a portfolio-grade backend project with Codex.

The target project is inspired by a small Go/Gin event REST API, but it is designed to show stronger production-minded engineering: explicit architecture, PostgreSQL constraints, secure auth, tests, OpenAPI documentation, Docker Compose, and CI.

## Current Implementation Status

Phase 2 database foundation is bootstrapped:

- Go module and API entrypoint.
- Explicit environment-based config loader.
- HTTP server with read, write, idle, and shutdown timeouts.
- `chi` router with request IDs and recovery middleware.
- `GET /healthz` process health endpoint.
- `GET /readyz` database readiness endpoint backed by PostgreSQL ping.
- Consistent JSON error envelope.
- PostgreSQL Docker Compose service.
- Reversible `golang-migrate` migration for users, refresh tokens, organizations, members, events, RSVPs, and audit logs.
- Database constraints for normalized emails, unique slugs, membership uniqueness, event state, RSVP uniqueness, enum-like status values, and audit metadata shape.
- Domain types plus repository interfaces and a Postgres store skeleton.
- Makefile targets for local development, migration commands, unit tests, vet, and integration migration tests.
- Initial OpenAPI spec for operational endpoints.

## Quick Start

Copy the example environment file if you want local overrides:

```bash
cp .env.example .env
```

Start PostgreSQL and apply migrations:

```bash
docker compose up -d postgres
make migrate-up
```

Run the API:

```bash
make run
```

Check the skeleton endpoints:

```bash
curl http://localhost:8080/healthz
curl http://localhost:8080/readyz
```

Run verification:

```bash
make test
make vet
make openapi-check
```

Run the Docker-backed migration smoke test when PostgreSQL is available:

```bash
make test-integration
```

Direct Go commands also work:

```bash
go test ./...
go vet ./...
```

## Target Project

GatherOps is a backend-focused event operations platform. It supports users, organizations, events, RSVP workflows, attendee management, audit logging, and operational readiness features such as health checks, migrations, OpenAPI docs, tests, Docker Compose, and CI.

The final implementation should show that the developer can build and reason about:

- API design and HTTP semantics.
- Authentication and authorization.
- Relational data modeling.
- Database migrations and constraints.
- Tests at unit, integration, and contract levels.
- Production-minded configuration and observability.
- Clear documentation and engineering tradeoffs.

## Configuration

| Variable | Default | Purpose |
| --- | --- | --- |
| `APP_ENV` | `development` | Application environment. Production requires an explicit `DATABASE_URL`. |
| `HTTP_PORT` | `8080` | API listen port. |
| `POSTGRES_PORT` | `5433` | Host port used by Docker Compose for local PostgreSQL. |
| `DATABASE_URL` | Local Docker Compose PostgreSQL URL on port `5433` | PostgreSQL connection string. |
| `HTTP_READ_TIMEOUT` | `5s` | Maximum duration for reading a request. |
| `HTTP_WRITE_TIMEOUT` | `10s` | Maximum duration before timing out response writes. |
| `HTTP_IDLE_TIMEOUT` | `60s` | Maximum idle keep-alive duration. |
| `SHUTDOWN_TIMEOUT` | `10s` | Graceful shutdown timeout. |
| `READINESS_TIMEOUT` | `2s` | Database readiness ping timeout. |
| `TEST_DATABASE_URL` | Local Docker Compose PostgreSQL URL on port `5433` | Optional database URL for integration tests. |

## API Shape

Operational endpoints are available at the root path:

| Method | Path | Purpose |
| --- | --- | --- |
| `GET` | `/healthz` | Returns `200` when the process is alive. |
| `GET` | `/readyz` | Returns `200` when PostgreSQL is reachable, otherwise `503`. |

Errors use the documented envelope:

```json
{
  "error": {
    "code": "database_unavailable",
    "message": "Database is not reachable",
    "requestId": "req_..."
  }
}
```

## Suggested Build Order

1. Bootstrap repository structure, tooling, Docker Compose, Makefile, config, and health endpoints.
2. Add database migrations and core schema.
3. Add users, auth, access tokens, refresh tokens, and password hashing.
4. Add organizations and membership roles.
5. Add events and owner/organizer authorization.
6. Add RSVP and attendee workflows with capacity and waitlist behavior.
7. Add audit logs, pagination, filtering, and OpenAPI docs.
8. Add tests, CI, observability, security hardening, and portfolio README polish.

## What To Preserve

Keep the planning files unless a deliberate ADR changes them. They are part of the portfolio story because they show engineering discipline, not just generated code.

## Known Limitations

- Product endpoints for auth, organizations, events, RSVP workflows, and audit logs are not implemented yet.
- Repository implementations currently stop at interfaces and a Postgres store skeleton; concrete query methods start in the auth slice.
- The OpenAPI spec currently covers only the operational endpoints.
- A hosted API docs UI route will be added in a later documentation slice.
