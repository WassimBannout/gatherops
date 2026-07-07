# Phase 1: Repository Foundation

## Goal

Create the initial code repository structure and a running API skeleton with health checks, config, local Postgres, and repeatable commands.

## Required Deliverables

- `go.mod`.
- `cmd/api/main.go`.
- Internal package layout.
- Config loader with environment variables.
- HTTP server with timeouts.
- Router with `/healthz` and `/readyz`.
- Consistent JSON error response helper.
- Docker Compose for PostgreSQL.
- `.env.example`.
- Makefile with at least `run`, `test`, `vet`, `migrate-up`, `migrate-down` placeholders.
- README quick start draft.

## Acceptance Criteria

- `go test ./...` passes.
- `go vet ./...` passes.
- `docker compose up -d` starts Postgres.
- `/healthz` returns 200.
- `/readyz` checks database connectivity.
- README explains how to start the skeleton.

## Codex Prompt Pointer

Use `prompts/00-repo-bootstrap.md`.


## Implementation Notes

- Router choice: `chi`, recorded in `decisions/ADR-0001-tech-stack.md`.
- Query strategy: not implemented in Phase 1; `sqlc` remains preferred for Phase 2 schema and repository work.
- Foundation scope: operational endpoints, config, server wiring, PostgreSQL connectivity, local tooling, and initial OpenAPI coverage.
- Known limitations: migration targets are placeholders until Phase 2, and the API docs UI route is not served yet.
