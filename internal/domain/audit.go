package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type AuditLog struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
	ActorUserID    uuid.UUID
	Action         string
	EntityType     string
	EntityID       uuid.UUID
	Metadata       json.RawMessage
	CreatedAt      time.Time
}
