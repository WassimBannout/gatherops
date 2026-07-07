# Codex Operating Instructions For GatherOps

You are building a portfolio-grade software engineering project, not a throwaway demo.

## Mission

Build GatherOps from scratch as a clean, tested, production-minded backend project. Use the planning documents in this repository as source-of-truth unless the human explicitly approves a change. If a requirement is ambiguous, make a conservative engineering assumption and document it.

## Non-Negotiables

- Do not copy code from the previous repository.
- Do not implement large unverified batches.
- Do not mark a task complete until tests and relevant docs are updated.
- Do not introduce hidden global state for core dependencies.
- Do not store plaintext secrets or tokens.
- Do not expose password hashes in API responses.
- Do not rely only on application checks when a database constraint is appropriate.
- Do not leave broken OpenAPI docs.

## Expected Engineering Style

- Prefer small vertical slices that include schema, handler/service/repository code, tests, and docs.
- Prefer explicit configuration over hard-coded production values.
- Keep business rules in services or domain functions, not scattered across handlers.
- Keep persistence details behind repository boundaries.
- Use request/response DTOs instead of exposing database models directly.
- Use structured error responses consistently.
- Use context propagation from HTTP request to database calls.
- Add indexes and constraints with migrations when behavior depends on them.

## Required Quality Gates

Before completing any implementation phase, run the relevant commands. The future implementation should eventually support:

```bash
go test ./...
go vet ./...
```

If a Makefile is added, prefer:

```bash
make test
make lint
make openapi-check
```

If a command cannot run because dependencies are not installed or Docker is unavailable, document the exact blocker and the command attempted.

## Documentation Discipline

When changing API behavior, update:

- OpenAPI spec.
- README examples.
- Any affected task checklist.
- Any relevant ADR if a design decision changes.

When adding a feature, include:

- What changed.
- How to run it locally.
- How it is tested.
- Known limitations.

## Security Discipline

Always consider:

- Authentication failure modes.
- Authorization boundaries.
- Token expiration and revocation.
- Input validation.
- Database constraints.
- Rate limiting for auth-sensitive endpoints.
- CORS and environment-specific configuration.

## AI-Native Collaboration Protocol

For each prompt from the human:

1. Read the relevant docs first.
2. State the implementation slice in one or two sentences.
3. Inspect the existing code before editing.
4. Make the smallest coherent change.
5. Run verification commands.
6. Summarize changed files, verification, and residual risks.

## Definition Of Done Pointer

Use `checklists/definition-of-done.md` for every phase.
