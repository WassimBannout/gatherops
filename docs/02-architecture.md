# Architecture

## Architectural Goal

Build a small but realistic backend service with clean boundaries, explicit dependencies, and testable business logic. Avoid over-engineering, but do not put all behavior directly in HTTP handlers.

## Proposed Runtime Architecture

```mermaid
flowchart TB
    Client[API client / Swagger UI] --> Router[HTTP router]
    Router --> Middleware[Middleware: request id, logging, CORS, auth]
    Middleware --> Handlers[HTTP handlers]
    Handlers --> Services[Application services]
    Services --> Policies[Authorization policies]
    Services --> Repos[Repository interfaces]
    Repos --> Postgres[(PostgreSQL)]
    Services --> TokenSvc[Token service]
    Services --> AuditSvc[Audit logger]
    Config[Environment config] --> Router
    Config --> Repos
    Config --> TokenSvc
```

## Suggested Go Package Layout

```text
.
|-- cmd
|   `-- api
|       `-- main.go
|-- internal
|   |-- app
|   |   `-- app.go
|   |-- config
|   |   `-- config.go
|   |-- httpapi
|   |   |-- router.go
|   |   |-- middleware.go
|   |   |-- handlers_auth.go
|   |   |-- handlers_orgs.go
|   |   |-- handlers_events.go
|   |   `-- errors.go
|   |-- domain
|   |   |-- users.go
|   |   |-- organizations.go
|   |   |-- events.go
|   |   `-- rsvps.go
|   |-- service
|   |   |-- auth_service.go
|   |   |-- org_service.go
|   |   |-- event_service.go
|   |   `-- rsvp_service.go
|   |-- repository
|   |   |-- interfaces.go
|   |   `-- postgres
|   |       `-- ...
|   |-- security
|   |   |-- password.go
|   |   `-- tokens.go
|   |-- observability
|   |   `-- logging.go
|   `-- validation
|       `-- validation.go
|-- migrations
|-- docs
|   `-- openapi.yaml
|-- test
|   `-- integration
|-- docker-compose.yml
|-- Makefile
|-- README.md
`-- AGENTS.md
```

## Dependency Direction

```mermaid
flowchart LR
    HTTP[httpapi] --> Service[service]
    Service --> Domain[domain]
    Service --> RepoInterfaces[repository interfaces]
    PostgresRepo[repository/postgres] --> RepoInterfaces
    PostgresRepo --> Domain
    HTTP --> DTOs[request/response DTOs]
    Security[security] --> Domain
    Service --> Security
```

Rules:

- `domain` should not import HTTP or database packages.
- `service` coordinates business logic and authorization checks.
- `httpapi` translates HTTP requests to service calls and service errors to HTTP responses.
- `repository/postgres` handles SQL and persistence details.
- `security` handles password hashing and token operations.

## Request Lifecycle

```mermaid
sequenceDiagram
    participant C as Client
    participant R as Router/Middleware
    participant H as Handler
    participant S as Service
    participant P as Policy
    participant DB as Repository/Postgres

    C->>R: HTTP request
    R->>R: request id, logging, auth if needed
    R->>H: route to handler
    H->>H: decode and validate DTO
    H->>S: call application service
    S->>P: authorize action
    S->>DB: read/write state
    DB-->>S: domain data
    S-->>H: result or typed error
    H-->>C: JSON response
```

## Error Handling

Use typed service errors and map them centrally to HTTP responses.

Suggested API error shape:

```json
{
  "error": {
    "code": "event_not_found",
    "message": "Event not found",
    "requestId": "req_123"
  }
}
```

## Configuration

All deployment-sensitive values should come from environment variables, with safe local defaults only for development.

Required eventual config:

- `APP_ENV`
- `HTTP_PORT`
- `DATABASE_URL`
- `JWT_ACCESS_SECRET`
- `JWT_REFRESH_SECRET` or refresh-token random generator config
- `ACCESS_TOKEN_TTL`
- `REFRESH_TOKEN_TTL`
- `CORS_ALLOWED_ORIGINS`
- `LOG_LEVEL`

## Why This Architecture Is Portfolio-Friendly

It is small enough to understand in one sitting but structured enough to show real engineering judgment. It avoids both extremes: a one-file demo and an overcomplicated enterprise skeleton.
