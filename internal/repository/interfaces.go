package repository

import (
	"context"

	"github.com/WassimBannout/gatherops/internal/domain"
	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, user domain.User) (domain.User, error)
	FindByID(ctx context.Context, id uuid.UUID) (domain.User, error)
	FindByEmail(ctx context.Context, email string) (domain.User, error)
}

type RefreshTokenRepository interface {
	Create(ctx context.Context, token domain.RefreshToken) (domain.RefreshToken, error)
	FindByHash(ctx context.Context, tokenHash string) (domain.RefreshToken, error)
	Revoke(ctx context.Context, id uuid.UUID) error
}

type OrganizationRepository interface {
	Create(ctx context.Context, organization domain.Organization) (domain.Organization, error)
	FindByID(ctx context.Context, id uuid.UUID) (domain.Organization, error)
	FindBySlug(ctx context.Context, slug string) (domain.Organization, error)
}

type OrganizationMemberRepository interface {
	Add(ctx context.Context, member domain.OrganizationMember) (domain.OrganizationMember, error)
	FindRole(ctx context.Context, organizationID, userID uuid.UUID) (domain.OrganizationRole, error)
}

type EventRepository interface {
	Create(ctx context.Context, event domain.Event) (domain.Event, error)
	FindByID(ctx context.Context, id uuid.UUID) (domain.Event, error)
}

type RSVPRepository interface {
	Upsert(ctx context.Context, rsvp domain.RSVP) (domain.RSVP, error)
	FindByEventAndUser(ctx context.Context, eventID, userID uuid.UUID) (domain.RSVP, error)
}

type AuditLogRepository interface {
	Append(ctx context.Context, entry domain.AuditLog) (domain.AuditLog, error)
}
