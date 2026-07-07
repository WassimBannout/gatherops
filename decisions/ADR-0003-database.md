# ADR-0003: Database And Data Access

## Status

Accepted

## Context

The previous project used SQLite and handwritten SQL. SQLite is fine for tutorials, but PostgreSQL better demonstrates production-style relational design.

## Decision

Use PostgreSQL with explicit migrations and database-backed constraints. Use `golang-migrate` for migration execution and keep persistence behind repository interfaces.

The Phase 2 foundation defines the schema for:

- `users`
- `refresh_tokens`
- `organizations`
- `organization_members`
- `events`
- `rsvps`
- `audit_logs`

The schema uses UUID primary keys, foreign keys, unique constraints, enum-like `CHECK` constraints, timestamp columns, and indexes for lookup/list paths described in `docs/04-data-model.md`.

Query implementation remains SQL-first. Feature slices should prefer `sqlc` for concrete query code when repository methods are implemented; carefully tested handwritten `pgx` remains acceptable only with explicit justification.

## Consequences

Positive:

- Real database constraints and transactions.
- Stronger portfolio signal.
- Better path for concurrent RSVP capacity behavior.
- Migration commands are repeatable through the Makefile.
- Integration tests can smoke-test migration up/down behavior against Docker PostgreSQL.

Tradeoffs:

- Requires Docker or local Postgres.
- Integration tests need database setup and are opt-in through `GATHEROPS_INTEGRATION_TESTS=1` or `make test-integration`.
- `golang-migrate` leaves its bookkeeping table after down migrations; tests account for this as migration metadata, not product schema.

## Required Schema Practices

- UUID primary keys.
- Explicit foreign keys.
- Unique constraints for natural uniqueness.
- Indexes for foreign keys and list filters.
- No `SELECT *` in application queries.
- Timestamps on important records.

## Implementation Notes

Phase 2 implemented the first reversible migration in `migrations/000001_create_core_schema.*.sql`, a small Go migration runner in `cmd/migrate`, and a migration smoke test in `test/integration`.
