# Security Requirements

## Security Goals

- Keep authentication predictable and testable.
- Avoid common JWT mistakes.
- Enforce authorization in service logic.
- Store secrets and tokens safely.
- Make dangerous defaults impossible in production.

## Authentication

Access tokens:

- Short-lived JWTs.
- Include standard `exp`, `iat`, `sub`, and token id where useful.
- Signed with strong environment-provided secret.
- Never accepted after expiration.

Refresh tokens:

- Long-lived opaque random tokens.
- Store hash only.
- Can be revoked on logout.
- Rotate on refresh if feasible.

Password hashing:

- Use bcrypt or argon2id.
- Never log passwords.
- Never return password hash in API responses.

## Authorization Model

Organization roles:

| Role | Permissions |
| --- | --- |
| owner | Manage organization, members, events, audit logs. |
| organizer | Manage events and attendees. |
| member | View organization resources and RSVP as normal user. |

Rules:

- Organization mutation requires owner.
- Event creation requires owner or organizer.
- Event update/publish/cancel requires owner or organizer in that organization.
- Attendee list requires owner or organizer.
- Public event read requires published status.
- Audit log read requires owner.

## Input Validation

Validate:

- Email format and normalization.
- Password minimum length.
- UUID path parameters.
- Date/time ranges.
- Capacity values.
- Pagination limits.
- Enum values.

## API Error Safety

- Do not leak internal SQL errors.
- Use generic login failure message.
- Include request id for debugging.
- Keep validation details useful but safe.

## CORS

- Default local origin can be permissive for development.
- Production must use explicit allowed origins.
- Do not combine wildcard origins with credentials.

## Rate Limiting

MVP should at least design for rate limiting. Implementation can be simple:

- Limit login attempts per IP/email.
- Limit register attempts per IP.
- Return `429` with consistent error shape.

## Security Tests

Required tests:

- Expired access token rejected.
- Token with missing `sub` rejected.
- Token signed with wrong secret rejected.
- Refresh token reuse after logout rejected.
- Member cannot perform owner actions.
- Organizer cannot manage organization owners.
- Non-member cannot mutate organization events.
- Password hash never appears in JSON responses.
