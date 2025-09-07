package gateway

import (
    "context"
    account "github.com/YukiOnishi1129/youtube-analytics/services/authority-service/internal/domain/account"
)

// RoleRepository manages role assignment for accounts.
type RoleRepository interface {
    // ListByAccount lists roles assigned to an account.
    ListByAccount(ctx context.Context, accountID string) ([]account.Role, error)

    // Assign assigns a role to an account.
    Assign(ctx context.Context, accountID string, role account.Role) error

    // Revoke revokes a role from an account.
    Revoke(ctx context.Context, accountID string, role account.Role) error
}

