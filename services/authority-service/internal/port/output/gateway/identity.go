package gateway

import (
    "context"
    account "github.com/YukiOnishi1129/youtube-analytics/services/authority-service/internal/domain/account"
)

// TokenVerifier verifies an ID token using an external identity provider
// and returns domain TokenClaims.
// Implementation example: Firebase Auth / OIDC JWKS.
type TokenVerifier interface {
    Verify(ctx context.Context, idToken string) (account.TokenClaims, error)
}

