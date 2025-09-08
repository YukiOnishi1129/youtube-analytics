package gateway

import (
    "context"
    "time"
    "github.com/YukiOnishi1129/youtube-analytics/services/authority-service/internal/domain"
)

// TokenVerifier verifies an ID token using an external identity provider.
// Implementation example: Firebase Auth / OIDC JWKS.
type TokenVerifier interface {
    Verify(ctx context.Context, idToken string) (claims TokenClaims, err error)
}

// TokenClaims is the verified claims required by authority use cases.
type TokenClaims struct {
    Subject       string
    Email         string
    EmailVerified bool
    DisplayName   string
    PhotoURL      string
    Provider      domain.Provider
    ProviderUID   string
}

// Clock abstracts time for use cases to enable deterministic tests.
type Clock interface { Now() time.Time }

// AuthTokens represents issued tokens from the identity provider.
type AuthTokens struct {
    IDToken      string
    RefreshToken string
    ExpiresIn    int64 // seconds
}

// IdentityProvider abstracts Identity Platform operations for MVP flows.
type IdentityProvider interface {
    // SignUp registers a new user and returns tokens.
    SignUp(ctx context.Context, email, password string) (AuthTokens, error)

    // SignIn authenticates a user and returns tokens.
    SignIn(ctx context.Context, email, password string) (AuthTokens, error)

    // SignOut revokes a refresh token (or all sessions).
    SignOut(ctx context.Context, refreshToken string) error

    // ResetPassword triggers out-of-band email for password reset.
    ResetPassword(ctx context.Context, email string) error
}

// ClaimsProvider provides current user's verified claims from context (via interceptor).
type ClaimsProvider interface {
    Current(ctx context.Context) (TokenClaims, error)
}
