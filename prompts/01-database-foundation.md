# Prompt 01: Database Foundation

```text
Implement the database foundation for GatherOps.

Read docs/04-data-model.md, decisions/ADR-0003-database.md, and tasks/00-master-plan.md. Add migrations for users, refresh_tokens, organizations, organization_members, events, rsvps, and audit_logs. Add required foreign keys, unique constraints, indexes, timestamps, and enum-like checks where appropriate.

Add migration commands to Makefile if not already present. Add integration or migration smoke tests if the test harness exists; otherwise add the minimal test harness needed.

Do not implement all business endpoints yet. This slice is schema, migration tooling, and repository boundaries/skeletons only.

Run verification commands and summarize results.
```
