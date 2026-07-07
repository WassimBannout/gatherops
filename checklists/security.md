# Security Checklist

## Auth

- [x] Access tokens include `exp`.
- [x] Expired access tokens are rejected.
- [x] Malformed tokens return 401 without panic.
- [x] Refresh tokens are stored hashed.
- [x] Logout revokes refresh tokens.
- [x] Password hash never appears in API output.

## Authorization

- [ ] Organization owner-only actions are protected.
- [ ] Organizer actions are protected.
- [ ] Non-members cannot mutate organization resources.
- [ ] Public reads expose only intended data.

## Input

- [ ] UUID path params are validated.
- [ ] Pagination limits are bounded.
- [ ] Date ranges are validated.
- [ ] Enum values are validated.
- [x] Duplicate user email returns 409.

## Database

- [x] Foreign keys exist for the Phase 2 core schema.
- [x] Unique constraints exist for membership and RSVP pairs.
- [ ] Transactions protect capacity-sensitive RSVP logic.
- [ ] No raw SQL string interpolation from user input.

## Ops

- [x] Production secrets have no unsafe defaults.
- [ ] CORS origins are configurable.
- [ ] Logs do not include passwords or tokens.
- [ ] Health/readiness endpoints do not leak secrets.
