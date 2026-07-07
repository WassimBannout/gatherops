# GatherOps

GatherOps is a production-style event operations API for organizations that need to create events, manage attendees, collect RSVPs, and audit operational changes. This repository contains the planning, specification, prompt workflow, and implementation workspace for building it as a portfolio-grade backend project with Codex.

The target project is inspired by a small Go/Gin event REST API, but it is designed to show stronger production-minded engineering: explicit architecture, PostgreSQL constraints, secure auth, tests, OpenAPI documentation, Docker Compose, and CI.

## Current Implementation Status

Phase 3 authentication and users are implemented:

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
- Domain types, repository interfaces, and concrete Postgres repositories for users and refresh tokens.
- `POST /api/v1/auth/register`, `POST /api/v1/auth/login`, `POST /api/v1/auth/refresh`, `POST /api/v1/auth/logout`, and `GET /api/v1/me`.
- Bcrypt password hashing, JWT access tokens with required `exp` and `sub`, opaque refresh tokens stored hashed at rest, refresh rotation, and logout revocation.
- Makefile targets for local development, migration commands, unit tests, vet, and integration migration/repository tests.
- OpenAPI spec for operational and authentication endpoints.

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

Check the operational endpoints:

```bash
curl http://localhost:8080/healthz
curl http://localhost:8080/readyz
```

Try the auth flow:

```bash
curl -s -X POST http://localhost:8080/api/v1/auth/register \
  -H 'Content-Type: application/json' \
  -d '{"name":"Ada Lovelace","email":"ada@example.com","password":"correct-password"}'
```

Use the returned `accessToken` as a bearer token for `GET /api/v1/me`, and use the returned `refreshToken` with `POST /api/v1/auth/refresh` or authenticated `POST /api/v1/auth/logout`.

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
| `JWT_ACCESS_SECRET` | Development-only local secret | HMAC secret for JWT access tokens. Required in production and must be at least 32 characters. |
| `ACCESS_TOKEN_TTL` | `15m` | Access token lifetime. |
| `REFRESH_TOKEN_TTL` | `720h` | Refresh token lifetime. |
| `TEST_DATABASE_URL` | Local Docker Compose PostgreSQL URL on port `5433` | Optional database URL for integration tests. |

## API Shape

Operational endpoints are available at the root path:

| Method | Path | Purpose |
| --- | --- | --- |
| `GET` | `/healthz` | Returns `200` when the process is alive. |
| `GET` | `/readyz` | Returns `200` when PostgreSQL is reachable, otherwise `503`. |

Authentication endpoints are available under `/api/v1`:

| Method | Path | Purpose |
| --- | --- | --- |
| `POST` | `/api/v1/auth/register` | Create an account and issue access/refresh tokens. |
| `POST` | `/api/v1/auth/login` | Exchange email/password for access/refresh tokens. |
| `POST` | `/api/v1/auth/refresh` | Rotate a refresh token and issue a new session. |
| `POST` | `/api/v1/auth/logout` | Revoke a refresh token; requires bearer access token. |
| `GET` | `/api/v1/me` | Return the authenticated user's profile. |

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

- Organization, event, RSVP, attendee, and audit-log product endpoints are not implemented yet.
- Rate limiting for auth-sensitive endpoints is not implemented yet.
- The API docs UI route is not served yet; the OpenAPI source lives at `docs/openapi.yaml`.
- Refresh tokens rotate on refresh and can be revoked on logout, but session listing and revoke-all-devices behavior are not implemented.
