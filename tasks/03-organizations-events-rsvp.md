# Phases 4-6: Organizations, Events, RSVP

## Goal

Implement the core product workflow: organizations create events, users RSVP, and organizers manage attendance.

## Organization Deliverables

- Organization schema and repository.
- Membership schema and repository.
- Organization service with owner creation.
- Member role changes.
- Authorization policy tests.

## Event Deliverables

- Event schema and repository.
- Event service.
- Draft, publish, update, cancel, and delete workflows.
- Public list endpoint with pagination.
- Public detail endpoint for published events.
- Organizer-only mutation endpoints.

## RSVP Deliverables

- RSVP schema and repository.
- RSVP service.
- Create/update/cancel RSVP.
- Capacity and waitlist behavior.
- Attendee list for organizers.

## Acceptance Criteria

- Only organization owners can manage members.
- Owners and organizers can manage events.
- Non-members cannot mutate organization events.
- Public users can see only published events.
- RSVP uniqueness is enforced in DB.
- Capacity cannot be exceeded by attending RSVPs.
- Waitlist behavior is deterministic and tested.
- All list endpoints are paginated.
