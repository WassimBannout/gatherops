# Phase 3: Authentication And Users

## Goal

Implement secure account creation and session management.

## Deliverables

- User migration and repository methods.
- Refresh token migration and repository methods.
- Password hashing helper.
- Access token generator and validator.
- Register endpoint.
- Login endpoint.
- Refresh endpoint.
- Logout endpoint.
- `GET /api/v1/me` endpoint.
- Auth middleware.
- OpenAPI definitions.

## Acceptance Criteria

- Duplicate emails return 409.
- Password hash never appears in responses.
- Expired access tokens return 401.
- Malformed tokens return 401 without panic.
- Logout revokes refresh token.
- Refresh token hash is stored, not raw token.
- Tests cover success and failure paths.

## Security Edge Cases

- Missing bearer prefix.
- Wrong signing secret.
- Missing `sub` claim.
- Expired `exp` claim.
- Reused revoked refresh token.
- Invalid password returns generic error.
