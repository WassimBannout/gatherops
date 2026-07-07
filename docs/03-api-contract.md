# API Contract

## API Principles

- JSON request and response bodies.
- Resource-oriented REST endpoints.
- Consistent error envelope.
- Pagination on all collection endpoints.
- Bearer access token for protected endpoints.
- OpenAPI 3.1 spec must stay current.

## Base Paths

```text
/api/v1
```

## Standard Error Shape

```json
{
  "error": {
    "code": "validation_failed",
    "message": "One or more fields are invalid",
    "details": {
      "email": "must be a valid email address"
    },
    "requestId": "req_01J..."
  }
}
```

`details` is optional.

## Auth Endpoints

| Method | Path | Auth | Purpose |
| --- | --- | --- | --- |
| POST | `/auth/register` | No | Create user account. |
| POST | `/auth/login` | No | Exchange credentials for tokens. |
| POST | `/auth/refresh` | No | Exchange refresh token for new access token. |
| POST | `/auth/logout` | Yes | Revoke current refresh token/session. |
| GET | `/me` | Yes | Return authenticated user's profile. |

## Organization Endpoints

| Method | Path | Auth | Purpose |
| --- | --- | --- | --- |
| POST | `/organizations` | Yes | Create organization. |
| GET | `/organizations` | Yes | List caller's organizations. |
| GET | `/organizations/{organizationId}` | Yes | Get organization detail. |
| POST | `/organizations/{organizationId}/members` | Owner | Add member by email. |
| PATCH | `/organizations/{organizationId}/members/{userId}` | Owner | Change member role. |
| DELETE | `/organizations/{organizationId}/members/{userId}` | Owner | Remove member. |

## Event Endpoints

| Method | Path | Auth | Purpose |
| --- | --- | --- | --- |
| POST | `/organizations/{organizationId}/events` | Organizer+ | Create draft event. |
| GET | `/events` | No | List published events with filters. |
| GET | `/events/{eventId}` | No | Get event detail. |
| PATCH | `/events/{eventId}` | Organizer+ | Update draft or published event. |
| POST | `/events/{eventId}/publish` | Organizer+ | Publish event. |
| POST | `/events/{eventId}/cancel` | Organizer+ | Cancel event. |
| DELETE | `/events/{eventId}` | Owner/Organizer | Delete draft event only. |

## RSVP Endpoints

| Method | Path | Auth | Purpose |
| --- | --- | --- | --- |
| PUT | `/events/{eventId}/rsvp` | Yes | Create or update caller RSVP. |
| DELETE | `/events/{eventId}/rsvp` | Yes | Cancel caller RSVP. |
| GET | `/events/{eventId}/attendees` | Organizer+ | List attendees and RSVP statuses. |

## Audit Endpoints

| Method | Path | Auth | Purpose |
| --- | --- | --- | --- |
| GET | `/organizations/{organizationId}/audit-logs` | Owner | List audit logs. |

## Operational Endpoints

| Method | Path | Auth | Purpose |
| --- | --- | --- | --- |
| GET | `/healthz` | No | Process is alive. |
| GET | `/readyz` | No | Process can reach dependencies. |
| GET | `/docs` | No | API documentation UI. |

## Pagination

Collection endpoints should accept:

```text
limit=20&offset=0
```

MVP can use offset pagination. Cursor pagination can be a later improvement.

Response shape:

```json
{
  "data": [],
  "pagination": {
    "limit": 20,
    "offset": 0,
    "total": 123
  }
}
```

## Event Filters

`GET /events` should support at minimum:

- `q`
- `from`
- `to`
- `organizationId`
- `status`

## HTTP Status Guidelines

| Status | Use |
| --- | --- |
| 200 | Successful read/update. |
| 201 | Successful creation. |
| 204 | Successful deletion or logout with no body. |
| 400 | Invalid syntax or invalid path/query parameter. |
| 401 | Missing, invalid, or expired authentication. |
| 403 | Authenticated but not authorized. |
| 404 | Resource not found or not visible to caller. |
| 409 | Uniqueness conflict or invalid state transition. |
| 422 | Semantically invalid payload if distinguished from 400. |
| 429 | Rate limit exceeded. |
| 500 | Unexpected server error. |
