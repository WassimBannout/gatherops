package domain

import (
	"time"

	"github.com/google/uuid"
)

type RSVPStatus string

const (
	RSVPStatusAttending  RSVPStatus = "attending"
	RSVPStatusWaitlisted RSVPStatus = "waitlisted"
	RSVPStatusDeclined   RSVPStatus = "declined"
	RSVPStatusCancelled  RSVPStatus = "cancelled"
)

type RSVP struct {
	ID        uuid.UUID
	EventID   uuid.UUID
	UserID    uuid.UUID
	Status    RSVPStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}
