# Master Implementation Plan

Build in phases. Each phase should end with a working, tested system.

## Phase 0: Planning Lock

- Read all docs and ADRs.
- Confirm stack decisions.
- Identify open questions.
- Create final implementation plan.

Done when Codex can explain the build plan and no major ambiguity remains.

## Phase 1: Repository Foundation

- Go module.
- Project layout.
- Config loading.
- HTTP server.
- Health/readiness endpoints.
- Structured errors.
- Docker Compose with PostgreSQL.
- Makefile.
- Initial README.

## Phase 2: Database Foundation

- Migration setup.
- Users, refresh tokens, organizations, members, events, RSVPs, audit logs schema.
- Repository interfaces and Postgres implementation skeleton.
- Migration tests.

## Phase 3: Authentication

- Register.
- Login.
- Refresh.
- Logout.
- Me endpoint.
- Password hashing.
- Access token middleware.
- Auth tests.

## Phase 4: Organizations

- Create organization.
- List organizations.
- Member management.
- Role policies.
- Authorization tests.

## Phase 5: Events

- Create draft event.
- Publish event.
- List public events.
- Event detail.
- Update/cancel/delete rules.
- Pagination and filters.

## Phase 6: RSVP And Attendees

- RSVP create/update/cancel.
- Capacity enforcement.
- Waitlist behavior.
- Organizer attendee list.
- Concurrency-sensitive tests.

## Phase 7: Audit, Docs, And Hardening

- Audit log writes for key actions.
- Audit log endpoint.
- OpenAPI completion.
- CI.
- Security checklist.
- Portfolio README polish.

## Phase 8: Optional Extensions

- Email notification outbox.
- Simple web dashboard.
- Metrics endpoint.
- Cursor pagination.
- Deployment manifest.


## Phase 2 Implementation Notes

Implemented in the current repository state:

- Reversible `golang-migrate` core schema migration.
- Makefile-backed `migrate-up`, `migrate-down`, and `test-integration` commands.
- Domain model types and repository interface boundaries.
- Postgres store skeleton for future concrete repositories.
- Docker-backed migration smoke test covering table creation, rollback, and representative constraints.

Remaining for later phases:

- Concrete repository query methods.
- Auth, organization, event, RSVP, and audit service behavior.
- OpenAPI coverage for product endpoints.
