# Prompt 02: Authentication Slice

```text
Implement the authentication and user slice from tasks/02-auth-users.md.

Requirements:
- Register.
- Login.
- Refresh.
- Logout.
- GET /api/v1/me.
- Password hashing.
- JWT access tokens with standard exp claim.
- Hashed refresh tokens.
- Auth middleware.
- Consistent error responses.
- OpenAPI updates.
- Tests for success, validation failure, duplicate email, expired token, malformed token, logout/revocation, and password hash not exposed.

Keep services testable and do not put all business rules in handlers.

Run tests, vet, and any OpenAPI validation command available.
```
