# Product Requirements

## Problem

Small teams need a lightweight way to publish events, manage RSVPs, enforce capacity, and understand who changed what. Existing tools are often too heavy for small organizations, while simple CRUD demos do not cover real operational needs.

## Goals

- Allow users to create accounts and authenticate securely.
- Allow users to create organizations and invite/manage members.
- Allow organization members to create and manage events.
- Allow users to RSVP to published events.
- Enforce event capacity and waitlist behavior.
- Provide audit history for important changes.
- Provide reliable API docs and local developer setup.

## Non-Goals For MVP

- Payment processing.
- Email delivery through a real provider.
- Native mobile apps.
- Complex calendar integrations.
- Multi-region deployment.
- Full admin dashboard.

## Current Implementation Status

As of July 7, 2026, GatherOps has completed the Phase 2 database foundation. The current implementation has a working API skeleton plus the core PostgreSQL schema, but it is not yet the full MVP product workflow.

Implemented:

- Go module, API entrypoint, and internal package layout.
- Environment-based configuration with safe local defaults and production database URL enforcement.
- HTTP server wiring with read, write, idle, and graceful shutdown timeouts.
- `chi` router with request IDs and recovery middleware.
- `GET /healthz` process health endpoint.
- `GET /readyz` database readiness endpoint backed by PostgreSQL ping.
- Consistent JSON error response envelope for readiness failures.
- Docker Compose PostgreSQL service for local development.
- `.env.example`, Makefile targets, and README quick start instructions.
- Reversible `golang-migrate` migration for the core schema.
- Database constraints and indexes for users, refresh tokens, organizations, members, events, RSVPs, and audit logs.
- Domain types, repository interfaces, and a Postgres store skeleton.
- Initial OpenAPI 3.1 file covering the operational endpoints.
- Focused tests for config loading, operational HTTP behavior, migration URL construction, and Docker-backed migration smoke coverage.

Requirement progress:

| Requirement | Current Status |
| --- | --- |
| PRD-001 | Schema support started. `users.email` is unique and normalized by database constraint, but registration is not implemented. |
| PRD-002 through PRD-004 | Schema support started for refresh tokens, but login, token issuing, hashing, refresh, and logout flows are not implemented. |
| PRD-005 through PRD-006 | Schema support started for organizations and membership roles, but organization APIs and role policies are not implemented. |
| PRD-007 through PRD-010 | Schema support started for events and RSVPs, including status, time-range, capacity, and uniqueness constraints. Event and RSVP APIs are not implemented. |
| PRD-011 | Schema support started for audit logs, but audit writes and audit APIs are not implemented. |
| PRD-012 | Partially started. An OpenAPI file exists for operational endpoints, but a hosted API docs route and full product contract are not implemented yet. |
| PRD-013 | Implemented for Phase 1. `/healthz` returns process health and `/readyz` checks PostgreSQL connectivity. |
| PRD-014 | Implemented for Phase 1. Docker Compose starts local PostgreSQL, defaulting to host port `5433`. |
| PRD-015 | Not started. CI will be added during the hardening/docs phase. |

Latest local verification for this foundation slice passed:

```bash
go test ./...
go vet ./...
make test
make lint
make openapi-check
make docker-up
make test-integration
make migrate-up
make migrate-down
```

## Personas

| Persona | Need | MVP Support |
| --- | --- | --- |
| Organization owner | Create organization, manage members, control events. | Yes |
| Event organizer | Create and publish events, view attendees. | Yes |
| Attendee | Discover events and RSVP. | Yes |
| API reviewer | Understand and run project quickly. | Yes |
| Operator | Know if API and DB are healthy. | Basic health/readiness |

## Core User Stories

### Authentication

- As a visitor, I can register with name, email, and password.
- As a user, I can log in and receive an access token and refresh token.
- As a user, I can refresh my session without re-entering my password.
- As a user, I can log out and revoke my refresh token.
- As a user, I can view my own profile.

### Organizations

- As a user, I can create an organization.
- As an organization owner, I can add members by email.
- As an organization owner, I can change member roles.
- As an organization member, I can view organizations I belong to.

### Events

- As an organizer, I can create a draft event.
- As an organizer, I can publish, update, cancel, or delete events according to business rules.
- As a public user, I can list published upcoming events.
- As a user, I can view event details.
- As an organizer, I can see attendee counts and RSVP statuses.

### RSVP

- As an authenticated user, I can RSVP to an event.
- As an authenticated user, I can change my RSVP status.
- As an authenticated user, I can cancel my RSVP.
- As a system, I should enforce capacity and place extra attendees on waitlist.

### Audit

- As an organization owner, I can see important changes made to organization events.
- As a reviewer, I can see that the system records meaningful operational events.

## MVP Functional Requirements

| ID | Requirement | Priority |
| --- | --- | --- |
| PRD-001 | Register users with unique normalized email. | Must |
| PRD-002 | Log in using email/password. | Must |
| PRD-003 | Issue short-lived access tokens. | Must |
| PRD-004 | Store refresh tokens hashed at rest. | Must |
| PRD-005 | Create organizations. | Must |
| PRD-006 | Support organization roles: owner, organizer, member. | Must |
| PRD-007 | Create, publish, update, cancel events. | Must |
| PRD-008 | List published events with pagination. | Must |
| PRD-009 | RSVP with attending, declined, waitlisted, cancelled statuses. | Must |
| PRD-010 | Enforce event capacity. | Must |
| PRD-011 | Record audit logs for auth, organization, event, and RSVP changes. | Should |
| PRD-012 | Serve OpenAPI docs. | Must |
| PRD-013 | Provide health and readiness endpoints. | Must |
| PRD-014 | Provide Docker Compose local environment. | Must |
| PRD-015 | Provide CI for tests and linting. | Should |

## MVP Acceptance Criteria

- All protected endpoints reject missing, expired, malformed, or invalid tokens.
- Users cannot mutate organizations they do not belong to.
- Members cannot perform owner-only actions.
- Events cannot exceed capacity without waitlisting.
- Duplicate RSVPs are prevented by database constraints.
- List endpoints use limit/offset or cursor pagination.
- API errors use a consistent JSON shape.
- Tests cover success, validation errors, authentication errors, and authorization errors.
