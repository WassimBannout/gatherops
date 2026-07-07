package domain

import (
	"time"

	"github.com/google/uuid"
)

type OrganizationRole string

const (
	OrganizationRoleOwner     OrganizationRole = "owner"
	OrganizationRoleOrganizer OrganizationRole = "organizer"
	OrganizationRoleMember    OrganizationRole = "member"
)

type Organization struct {
	ID        uuid.UUID
	Name      string
	Slug      string
	CreatedBy uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
}

type OrganizationMember struct {
	OrganizationID uuid.UUID
	UserID         uuid.UUID
	Role           OrganizationRole
	CreatedAt      time.Time
}
