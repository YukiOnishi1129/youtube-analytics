package domain

import (
	"net"
	"time"

	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain/valueobject"
)

// AuditLog represents an audit trail entry for administrative actions
type AuditLog struct {
	ID           valueobject.UUID
	ActorID      valueobject.UUID
	ActorEmail   string
	Action       string
	ResourceType string
	ResourceID   string
	OldValues    map[string]interface{}
	NewValues    map[string]interface{}
	IPAddress    net.IP
	UserAgent    string
	CreatedAt    time.Time
}

// NewAuditLog creates a new audit log entry
func NewAuditLog(
	id valueobject.UUID,
	actorID valueobject.UUID,
	actorEmail string,
	action string,
	resourceType string,
	resourceID string,
	oldValues map[string]interface{},
	newValues map[string]interface{},
	ipAddress net.IP,
	userAgent string,
) (*AuditLog, error) {
	if actorID == "" {
		return nil, ErrInvalidInput
	}
	if actorEmail == "" {
		return nil, ErrInvalidInput
	}
	if action == "" {
		return nil, ErrInvalidInput
	}
	if resourceType == "" {
		return nil, ErrInvalidInput
	}

	return &AuditLog{
		ID:           id,
		ActorID:      actorID,
		ActorEmail:   actorEmail,
		Action:       action,
		ResourceType: resourceType,
		ResourceID:   resourceID,
		OldValues:    oldValues,
		NewValues:    newValues,
		IPAddress:    ipAddress,
		UserAgent:    userAgent,
		CreatedAt:    time.Now(),
	}, nil
}