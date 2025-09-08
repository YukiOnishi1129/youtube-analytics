package security

import (
    "context"
    outgateway "github.com/YukiOnishi1129/youtube-analytics/services/authority-service/internal/port/output/gateway"
)

type ctxKey string

const claimsKey ctxKey = "auth.claims"

// ContextClaimsProvider implements ClaimsProvider using context values.
type ContextClaimsProvider struct{}

func (p *ContextClaimsProvider) Current(ctx context.Context) (outgateway.TokenClaims, error) {
    v := ctx.Value(claimsKey)
    if v == nil {
        return outgateway.TokenClaims{}, ErrNoClaims
    }
    c, ok := v.(outgateway.TokenClaims)
    if !ok {
        return outgateway.TokenClaims{}, ErrNoClaims
    }
    return c, nil
}

var ErrNoClaims = ErrUnauthorized("missing claims in context")

type ErrUnauthorized string

func (e ErrUnauthorized) Error() string { return string(e) }

// WithClaims attaches claims into context.
func WithClaims(ctx context.Context, claims outgateway.TokenClaims) context.Context {
    return context.WithValue(ctx, claimsKey, claims)
}

// Ensure ContextClaimsProvider implements the output port.
var _ outgateway.ClaimsProvider = (*ContextClaimsProvider)(nil)
