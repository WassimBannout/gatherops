# ADR-0002: Authentication Strategy

## Status

Accepted

## Context

The previous project used JWT access tokens but had a critical typo in the expiry claim. A portfolio project should show secure token lifecycle decisions.

## Decision

Use short-lived JWT access tokens plus opaque refresh tokens stored hashed in the database.

Access tokens:

- JWT signed with HMAC SHA-256.
- Include `sub`, `exp`, `iat`, and `jti` claims.
- Default local TTL is 15 minutes through `ACCESS_TOKEN_TTL`.
- `JWT_ACCESS_SECRET` is required in production and must be at least 32 characters.

Refresh tokens:

- Cryptographically random opaque value.
- Stored as SHA-256 hash only.
- Default local TTL is 30 days through `REFRESH_TOKEN_TTL`.
- Rotated on refresh.
- Revoked on logout.

Passwords:

- Hashed with bcrypt.
- Never returned in API responses.

## Consequences

Positive:

- Demonstrates realistic session management.
- Allows logout/revocation for refresh sessions.
- Avoids indefinite JWT validity.
- Keeps refresh token compromise impact lower by storing only hashes.

Tradeoffs:

- More implementation work than access-token-only auth.
- Requires careful tests around refresh rotation and revocation.
- Session listing and revoke-all-devices are out of scope for the current slice.

## Required Tests

- Expired access token rejected.
- Malformed access token rejected.
- Token with missing `sub` rejected.
- Token signed with wrong secret rejected.
- Refresh token can create a new session.
- Revoked refresh token cannot be reused.
- Password hash is never exposed.

## Implementation Notes

Phase 3 implemented register, login, refresh, logout, `GET /api/v1/me`, bcrypt password hashing, JWT access tokens, hashed refresh tokens, auth middleware, Postgres user/session repositories, and OpenAPI coverage for auth endpoints.
