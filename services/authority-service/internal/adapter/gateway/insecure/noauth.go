package insecure

import (
	"context"

	outgateway "github.com/YukiOnishi1129/youtube-analytics/services/authority-service/internal/port/output/gateway"
)

// NoopVerifier accepts any token and returns fixed claims for local dev.
type NoopVerifier struct{}

func (NoopVerifier) Verify(ctx context.Context, idToken string) (outgateway.TokenClaims, error) {
	return outgateway.TokenClaims{
		Subject:       "dev-user",
		Email:         "dev@example.com",
		EmailVerified: true,
		DisplayName:   "Dev User",
		PhotoURL:      "",
	}, nil
}

// DummyIDP is a local stub for IdentityProvider operations.
type DummyIDP struct{}

func (DummyIDP) SignUp(ctx context.Context, email, password string) (outgateway.AuthTokens, error) {
	return outgateway.AuthTokens{IDToken: "dev-id-token", RefreshToken: "dev-refresh-token", ExpiresIn: 3600}, nil
}

func (DummyIDP) SignIn(ctx context.Context, email, password string) (outgateway.AuthTokens, error) {
	return outgateway.AuthTokens{IDToken: "dev-id-token", RefreshToken: "dev-refresh-token", ExpiresIn: 3600}, nil
}

func (DummyIDP) SignOut(ctx context.Context, refreshToken string) error { return nil }

func (DummyIDP) ResetPassword(ctx context.Context, email string) error { return nil }
