package firebase

import (
	"context"

	outgateway "github.com/YukiOnishi1129/youtube-analytics/services/authority-service/internal/port/output/gateway"
	"github.com/coreos/go-oidc/v3/oidc"
)

type OIDCVerifier struct {
	v *oidc.IDTokenVerifier
}

// NewOIDCVerifier constructs a TokenVerifier for Firebase (or any OIDC provider).
// issuer example: https://securetoken.google.com/<PROJECT_ID>
func NewOIDCVerifier(ctx context.Context, issuer string, audience string) (*OIDCVerifier, error) {
	provider, err := oidc.NewProvider(ctx, issuer)
	if err != nil {
		return nil, err
	}
	config := &oidc.Config{ClientID: audience}
	return &OIDCVerifier{v: provider.Verifier(config)}, nil
}

var _ outgateway.TokenVerifier = (*OIDCVerifier)(nil)

func (o *OIDCVerifier) Verify(ctx context.Context, idToken string) (outgateway.TokenClaims, error) {
	tok, err := o.v.Verify(ctx, idToken)
	if err != nil {
		return outgateway.TokenClaims{}, err
	}
	var claims struct {
		Sub           string `json:"sub"`
		Email         string `json:"email"`
		EmailVerified bool   `json:"email_verified"`
		Name          string `json:"name"`
		Picture       string `json:"picture"`
	}
	if err := tok.Claims(&claims); err != nil {
		return outgateway.TokenClaims{}, err
	}
	return outgateway.TokenClaims{
		Subject:       claims.Sub,
		Email:         claims.Email,
		EmailVerified: claims.EmailVerified,
		DisplayName:   claims.Name,
		PhotoURL:      claims.Picture,
	}, nil
}
