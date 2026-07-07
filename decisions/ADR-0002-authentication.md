# ADR-0002: Authentication Strategy

## Status

Proposed

## Context

The previous project used JWT access tokens but had a critical typo in the expiry claim. A portfolio project should show secure token lifecycle decisions.

## Decision

Use short-lived JWT access tokens plus opaque refresh tokens stored hashed in the database.

Access tokens:

- JWT.
- Include `sub`, `exp`, `iat`.
- Short TTL, for example 15 minutes.

Refresh tokens:

- Cryptographically random opaque value.
- Stored as hash only.
- Longer TTL, for example 7 to 30 days.
- Revoked on logout.

## Consequences

Positive:

- Demonstrates realistic session management.
- Allows logout/revocation for refresh sessions.
- Avoids indefinite JWT validity.

Tradeoffs:

- More implementation work than access-token-only auth.
- Requires careful tests around refresh and revocation.

## Required Tests

- Expired access token rejected.
- Malformed access token rejected.
- Refresh token can create a new access token.
- Revoked refresh token cannot be reused.
- Password hash is never exposed.
