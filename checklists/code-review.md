# Code Review Checklist

## Architecture

- Does the change belong in handler, service, domain, repository, or security package?
- Are business rules testable without HTTP?
- Is persistence hidden behind a boundary?
- Does the change avoid cyclic dependencies?

## API

- Are request DTOs separate from database models?
- Are response fields intentional?
- Are errors consistent?
- Are status codes correct?
- Is OpenAPI updated?

## Database

- Are constraints enforced in the database when needed?
- Are indexes added for new query patterns?
- Are migrations reversible when practical?
- Are queries explicit and not `SELECT *`?
- Are transactions used for multi-step state changes?

## Security

- Is auth required where needed?
- Are roles checked centrally and tested?
- Are tokens handled safely?
- Are secrets configurable?
- Could error messages leak sensitive data?

## Tests

- Are failure paths covered?
- Are authorization boundaries covered?
- Are database constraints tested?
- Does the test data make intent clear?
