package postgres

import (
	"context"

	"github.com/YukiOnishi1129/youtube-analytics/services/authority-service/internal/adapter/gateway/postgres/sqlcgen"
	"github.com/YukiOnishi1129/youtube-analytics/services/authority-service/internal/domain"
	"github.com/google/uuid"
)

// RoleRepository methods
func (r *Repository) ListRolesByAccount(ctx context.Context, accountID string) ([]domain.Role, error) {
	acc, err := uuid.Parse(accountID)
	if err != nil {
		return nil, nil
	}
	rows, err := r.q.ListRolesByAccount(ctx, acc)
	if err != nil {
		return nil, err
	}
	roles := make([]domain.Role, 0, len(rows))
	for _, name := range rows {
		roles = append(roles, domain.Role{Name: name})
	}
	return roles, nil
}

func (r *Repository) Assign(ctx context.Context, accountID string, role domain.Role) error {
	acc, err := uuid.Parse(accountID)
	if err != nil {
		return err
	}
	return r.q.AssignRole(ctx, sqlcgen.AssignRoleParams{AccountID: acc, Name: role.Name})
}

func (r *Repository) Revoke(ctx context.Context, accountID string, role domain.Role) error {
	acc, err := uuid.Parse(accountID)
	if err != nil {
		return err
	}
	return r.q.RevokeRole(ctx, sqlcgen.RevokeRoleParams{AccountID: acc, Name: role.Name})
}
