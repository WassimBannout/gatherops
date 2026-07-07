# DevOps And Observability

## Local Development Goals

A reviewer should be able to run the project locally without manual database setup.

Expected local stack:

- API service running on local port.
- PostgreSQL through Docker Compose.
- Migrations command.
- OpenAPI docs route.

## Recommended Files For Final Implementation

```text
docker-compose.yml
Dockerfile
Makefile
.env.example
.github/workflows/ci.yml
```

## Health Endpoints

### `GET /healthz`

Returns 200 if the process is alive.

Example:

```json
{
  "status": "ok"
}
```

### `GET /readyz`

Returns 200 only if dependencies are reachable, especially PostgreSQL.

Example:

```json
{
  "status": "ready",
  "dependencies": {
    "database": "ok"
  }
}
```

## Logging

Use structured logs with at least:

- Timestamp.
- Level.
- Request id.
- Method.
- Path.
- Status.
- Latency.
- Error code if present.

## Request IDs

Every request should get or preserve a request id. Include it in logs and error responses.

## Metrics

Metrics are optional for MVP, but the design should allow adding:

- Request count.
- Request latency.
- Error count.
- Database query latency.

## Deployment Story

The portfolio README should explain:

- How to run locally.
- How configuration works.
- How migrations are applied.
- What would change for production.
- Why PostgreSQL is used.
- Known limitations.
