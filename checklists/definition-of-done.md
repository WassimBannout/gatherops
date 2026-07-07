# Definition Of Done

A task is done only when all relevant items are true.

## Functionality

- The feature satisfies the stated acceptance criteria.
- Happy path and important failure paths work.
- Edge cases are handled deliberately.
- API status codes match `docs/03-api-contract.md`.

## Code Quality

- Code is formatted.
- Naming is clear.
- Dependencies flow in the intended direction.
- No unrelated refactors are included.
- No production TODOs remain unless documented as known limitations.

## Tests

- Unit tests added where business logic exists.
- Handler tests added for HTTP behavior.
- Integration tests added for database behavior.
- Existing tests pass.

## Security

- Authentication and authorization paths are tested.
- Secrets are not hard-coded for production.
- Passwords/tokens are not logged or returned.
- Database constraints back important invariants.

## Documentation

- OpenAPI updated if API changed.
- README updated if setup or behavior changed.
- ADR updated if architecture changed.
- Example commands remain valid.

## Verification

Run and record results:

```bash
go test ./...
go vet ./...
```

Use Make targets once available.
