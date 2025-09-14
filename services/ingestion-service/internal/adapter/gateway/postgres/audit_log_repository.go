package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"net"

	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/adapter/gateway/postgres/sqlcgen"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain/valueobject"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/port/output/gateway"
	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
)

// auditLogRepository implements gateway.AuditLogRepository interface
type auditLogRepository struct {
	*Repository
}

// NewAuditLogRepository creates a new audit log repository
func NewAuditLogRepository(repo *Repository) gateway.AuditLogRepository {
	return &auditLogRepository{Repository: repo}
}

// Save creates a new audit log entry
func (r *auditLogRepository) Save(ctx context.Context, log *domain.AuditLog) error {
	id, err := uuid.Parse(string(log.ID))
	if err != nil {
		return err
	}

	actorID, err := uuid.Parse(string(log.ActorID))
	if err != nil {
		return err
	}

	var oldValues, newValues pqtype.NullRawMessage
	if log.OldValues != nil {
		data, err := json.Marshal(log.OldValues)
		if err != nil {
			return err
		}
		oldValues = pqtype.NullRawMessage{RawMessage: data, Valid: true}
	}

	if log.NewValues != nil {
		data, err := json.Marshal(log.NewValues)
		if err != nil {
			return err
		}
		newValues = pqtype.NullRawMessage{RawMessage: data, Valid: true}
	}

	var ipAddress pqtype.Inet
	if log.IPAddress != nil {
		ipAddress = pqtype.Inet{IPNet: &net.IPNet{IP: log.IPAddress}, Valid: true}
	}

	return r.q.CreateAuditLog(ctx, sqlcgen.CreateAuditLogParams{
		ID:           id,
		ActorID:      actorID,
		ActorEmail:   log.ActorEmail,
		Action:       log.Action,
		ResourceType: log.ResourceType,
		ResourceID:   toNullString(&log.ResourceID),
		OldValues:    oldValues,
		NewValues:    newValues,
		IpAddress:    ipAddress,
		UserAgent:    toNullString(&log.UserAgent),
		CreatedAt:    log.CreatedAt,
	})
}

// FindByActor finds audit logs by actor ID
func (r *auditLogRepository) FindByActor(ctx context.Context, actorID valueobject.UUID, limit, offset int) ([]*domain.AuditLog, error) {
	aid, err := uuid.Parse(string(actorID))
	if err != nil {
		return nil, err
	}

	rows, err := r.q.ListAuditLogsByActor(ctx, sqlcgen.ListAuditLogsByActorParams{
		ActorID: aid,
		Limit:   int32(limit),
		Offset:  int32(offset),
	})
	if err != nil {
		return nil, err
	}

	logs := make([]*domain.AuditLog, len(rows))
	for i, row := range rows {
		logs[i] = toDomainAuditLog(row)
	}

	return logs, nil
}

// FindByResource finds audit logs by resource type and ID
func (r *auditLogRepository) FindByResource(ctx context.Context, resourceType, resourceID string, limit, offset int) ([]*domain.AuditLog, error) {
	rows, err := r.q.ListAuditLogsByResource(ctx, sqlcgen.ListAuditLogsByResourceParams{
		ResourceType: resourceType,
		ResourceID:   sql.NullString{String: resourceID, Valid: resourceID != ""},
		Limit:        int32(limit),
		Offset:       int32(offset),
	})
	if err != nil {
		return nil, err
	}

	logs := make([]*domain.AuditLog, len(rows))
	for i, row := range rows {
		logs[i] = toDomainAuditLog(row)
	}

	return logs, nil
}

// FindRecent finds recent audit logs
func (r *auditLogRepository) FindRecent(ctx context.Context, limit int) ([]*domain.AuditLog, error) {
	rows, err := r.q.ListRecentAuditLogs(ctx, int32(limit))
	if err != nil {
		return nil, err
	}

	logs := make([]*domain.AuditLog, len(rows))
	for i, row := range rows {
		logs[i] = toDomainAuditLog(row)
	}

	return logs, nil
}

// toDomainAuditLog converts a database row to a domain audit log
func toDomainAuditLog(row sqlcgen.IngestionAuditLog) *domain.AuditLog {
	log := &domain.AuditLog{
		ID:           valueobject.UUID(row.ID.String()),
		ActorID:      valueobject.UUID(row.ActorID.String()),
		ActorEmail:   row.ActorEmail,
		Action:       row.Action,
		ResourceType: row.ResourceType,
		ResourceID:   fromNullString(row.ResourceID),
		UserAgent:    fromNullString(row.UserAgent),
		CreatedAt:    row.CreatedAt,
	}

	if row.OldValues.Valid {
		var oldValues map[string]interface{}
		if err := json.Unmarshal(row.OldValues.RawMessage, &oldValues); err == nil {
			log.OldValues = oldValues
		}
	}

	if row.NewValues.Valid {
		var newValues map[string]interface{}
		if err := json.Unmarshal(row.NewValues.RawMessage, &newValues); err == nil {
			log.NewValues = newValues
		}
	}

	if row.IpAddress.Valid && row.IpAddress.IPNet != nil {
		log.IPAddress = row.IpAddress.IPNet.IP
	}

	return log
}

// fromNullString converts sql.NullString to string
func fromNullString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}