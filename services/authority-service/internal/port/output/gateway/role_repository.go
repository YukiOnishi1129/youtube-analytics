package gateway

import (
    "context"
    "github.com/YukiOnishi1129/youtube-analytics/services/authority-service/internal/domain"
)

// RoleRepository manages role assignment for accounts.
type RoleRepository interface {
    // ListRolesByAccount lists roles assigned to an account.
    ListRolesByAccount(ctx context.Context, accountID string) ([]domain.Role, error)
    // Assign assigns a role to an account.
    Assign(ctx context.Context, accountID string, role domain.Role) error
    // Revoke revokes a role from an account.
    Revoke(ctx context.Context, accountID string, role domain.Role) error
}

