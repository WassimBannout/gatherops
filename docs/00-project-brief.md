# Project Brief

## Name

GatherOps

## One-Sentence Pitch

GatherOps is a production-style event operations API that helps organizations create events, manage attendees, collect RSVPs, and audit operational changes.

## Why This Is A Strong Portfolio Project

A basic CRUD API is usually not enough to stand out. GatherOps is intentionally scoped to demonstrate the skills employers look for in software engineering candidates:

- Modeling a real domain with users, organizations, roles, events, RSVPs, and audit logs.
- Designing secure authentication and authorization flows.
- Building reliable persistence with migrations, indexes, and constraints.
- Writing tests beyond happy paths.
- Documenting API contracts and tradeoffs.
- Running the system locally with predictable developer tooling.
- Explaining engineering choices through ADRs and a polished README.

## Inspiration From The Existing Project

The previous project already had useful foundations:

- Go backend.
- Gin HTTP routes.
- User registration and login.
- JWT auth.
- Event CRUD.
- Attendee relationship table.
- SQLite persistence.
- SQL migrations.
- Swagger docs.

GatherOps keeps the educational spirit but improves the engineering bar:

- PostgreSQL instead of local-only SQLite.
- Stronger schema constraints.
- Access and refresh token lifecycle.
- Clear role-based authorization.
- Request/response DTOs.
- OpenAPI-first documentation.
- Tests from the start.
- Docker Compose and CI.
- Health/readiness endpoints.
- Portfolio-ready documentation.

## Target Audience

The product audience is small organizations that host events: meetups, clubs, internal company groups, university societies, bootcamps, and community teams.

The portfolio audience is hiring managers and engineers reviewing the candidate's GitHub profile.

## Scope Boundary

The first version is backend-only. A small frontend can be added later, but the backend should stand on its own.

## Success Criteria

A reviewer should be able to:

1. Clone the repository.
2. Run one command to start dependencies.
3. Run migrations.
4. Start the API.
5. Open the API docs.
6. Register/login.
7. Create an organization and event.
8. RSVP to an event.
9. Run the test suite.
10. Understand the architecture from README and diagrams.
