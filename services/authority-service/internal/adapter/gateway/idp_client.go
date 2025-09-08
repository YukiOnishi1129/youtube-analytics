package gateway

import (
    "context"
    outgateway "github.com/YukiOnishi1129/youtube-analytics/services/authority-service/internal/port/output/gateway"
)

// IDPClient is a stub implementation of IdentityProvider.
// Replace with Firebase/Identity Platform client.
type IDPClient struct{}

var _ outgateway.IdentityProvider = (*IDPClient)(nil)

func (c *IDPClient) SignUp(ctx context.Context, email, password string) (outgateway.AuthTokens, error) {
    return outgateway.AuthTokens{IDToken: "stub", RefreshToken: "stub", ExpiresIn: 3600}, nil
}

func (c *IDPClient) SignIn(ctx context.Context, email, password string) (outgateway.AuthTokens, error) {
    return outgateway.AuthTokens{IDToken: "stub", RefreshToken: "stub", ExpiresIn: 3600}, nil
}

func (c *IDPClient) SignOut(ctx context.Context, refreshToken string) error { return nil }

func (c *IDPClient) ResetPassword(ctx context.Context, email string) error { return nil }

