# Prompt 00: Bootstrap The Implementation Repository

Use this prompt when you are ready for Codex to start creating the actual application code in this repository.

```text
Read AGENTS.md, project.yaml, docs/*.md, decisions/*.md, tasks/*.md, and checklists/*.md.

We are building GatherOps from scratch as a portfolio-grade backend project. Do not copy code from the previous repository. Use the docs in this repo as source-of-truth.

First, summarize the target architecture and confirm any stack choices that are still proposed, especially router choice and query strategy. Then implement Phase 1 from tasks/01-foundation.md:

- Go module and project layout.
- Config loader.
- HTTP server with timeouts.
- Router with /healthz and /readyz.
- Consistent JSON error response helper.
- Docker Compose for PostgreSQL.
- .env.example.
- Makefile.
- README quick start draft.

Keep the slice small and verifiable. Run go test ./... and go vet ./.... Summarize changed files, verification results, and any blockers.
```
