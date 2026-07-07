package domain

import (
	"time"

	"github.com/google/uuid"
)

type EventStatus string

const (
	EventStatusDraft     EventStatus = "draft"
	EventStatusPublished EventStatus = "published"
	EventStatusCancelled EventStatus = "cancelled"
)

type Event struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
	CreatedBy      uuid.UUID
	Title          string
	Description    string
	Location       string
	StartsAt       time.Time
	EndsAt         time.Time
	Capacity       *int
	Status         EventStatus
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
