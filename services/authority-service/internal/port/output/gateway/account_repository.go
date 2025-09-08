package gateway

import (
    "context"
    "github.com/YukiOnishi1129/youtube-analytics/services/authority-service/internal/domain"
)

// AccountRepository owns persistence operations for Account aggregate.
type AccountRepository interface {
    // FindByID returns the account or ErrNotFound.
    FindByID(ctx context.Context, id string) (*domain.Account, error)
    // FindByEmail returns the account by email (case-insensitive) or ErrNotFound.
    FindByEmail(ctx context.Context, email string) (*domain.Account, error)
    // Save creates or updates an account.
    Save(ctx context.Context, a *domain.Account) error
}

